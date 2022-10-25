package postgresql

import (
	"context"
	"encoding/csv"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"golang.org/x/sync/errgroup"
	"io"
	"os"
	"postgres-perf/configuration"
	"postgres-perf/utils/logger"
	"strings"
	"time"
)

// Provider hold postgres related variables
type Provider struct {
	PostgresqlHost     string
	PostgresqlPort     int32
	PostgresqlUser     string
	PostgreslDatabase  string
	PostgresqlPassword string
	Workers            int32
}

type QueryParams struct {
	Host  string
	Start string
	End   string
}

var (
	InvalidHeader = fmt.Errorf("headers missing")
)

type DbCon interface {
	Prepare(ctx context.Context, name string, sql string) (*pgconn.StatementDescription, error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	Close(ctx context.Context) error
}

type RunQuery struct {
	Connection func(context.Context) (DbCon, error)
	Query      func(context.Context, DbCon, []QueryParams) ([]Stat, error)
}

type ResponseVerificationProperties struct {
	Message string `json:"@odata.context"`
}

func New(ctx context.Context) (*Provider, error) {
	var (
		//err       error
		provider *Provider
	)

	log := logger.SugaredLogger().WithContextCorrelationId(ctx)
	log.Debug("Instantiating Postgres provider")

	appConfig := configuration.AppConfig()
	provider = &Provider{}

	// check if pgsql server set in env var
	if appConfig.PostgresqlHost == "" {
		return nil, fmt.Errorf("pgsql endpoint not set")
	} else {
		provider.PostgresqlHost = appConfig.PostgresqlHost
		log.Debugf("POSTGRESQL_ENDPOINTS:%s", provider.PostgresqlHost)
	}

	// check if pgsql server set in env var
	if appConfig.PostgresqlPort == 0 {
		return nil, fmt.Errorf("pgsql port not set")
	} else {
		provider.PostgresqlPort = appConfig.PostgresqlPort
		log.Debugf("POSTGRESQL_PORT:%v", provider.PostgresqlPort)
	}

	// check workers
	if appConfig.Workers == 0 {
		return nil, fmt.Errorf("workers not set")
	} else {
		provider.Workers = appConfig.Workers
		log.Debugf("WORKERS:%v", provider.Workers)
	}

	// check if pgsql user is set in env var
	if appConfig.PostgresqlUser == "" {
		return nil, fmt.Errorf("pgsql user not set")
	} else {
		provider.PostgresqlUser = appConfig.PostgresqlUser
		log.Debugf("POSTGRESQL_USER:%s", provider.PostgresqlUser)
	}

	// check if pgsql db is set in env var
	if appConfig.PostgresqlDatabase == "" {
		return nil, fmt.Errorf("pgsql db not set")
	} else {
		provider.PostgreslDatabase = appConfig.PostgresqlDatabase
		log.Debugf("POSTGRESQL_DATABASE:%s", provider.PostgreslDatabase)
	}

	// check if db passwd is set in env var
	if appConfig.PostgresqlPassword == "" {
		return nil, fmt.Errorf("pgsql password not set")
	} else {
		provider.PostgresqlPassword = appConfig.PostgresqlPassword
		log.Debugf("POSTGRESQL_PASSWORD:%s", provider.PostgresqlPassword)
	}
	return provider, nil
}

// ConnectionString connection string
func (p *Provider) ConnectionString(ctx context.Context) string {
	log := logger.SugaredLogger().WithContextCorrelationId(ctx).With("package", "pg", "action", "ConnectionString")
	log.Debugf("init db connection string")
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s", p.PostgresqlUser, p.PostgresqlPassword,
		p.PostgresqlHost, p.PostgresqlPort, p.PostgreslDatabase)
}

// QueryFileParams open file
func (p *Provider) QueryFileParams(ctx context.Context, File string) (io.ReadCloser, error) {
	log := logger.SugaredLogger().WithContextCorrelationId(ctx).With("package", "pg", "action", "QueryFileParams")
	log.Debugf("query file param")
	if File == "" {
		return io.NopCloser(os.Stdin), nil
	} else {
		return os.Open(File)
	}
}

// ValidateHeader check if headers are in cvs file
func (p *Provider) ValidateHeader(line []string) error {
	HeaderStr := strings.Join(line, ",")
	if HeaderStr == configuration.ExpectedHeaders {
		return nil
	}
	return fmt.Errorf("expected: %s but give: %s, %w", configuration.ExpectedHeaders, HeaderStr, InvalidHeader)
}

// LoadCSV load content
func (p *Provider) LoadCSV(reader io.Reader) ([]QueryParams, error) {
	csvreader := csv.NewReader(reader)

	var keys = map[string]int{}
	read, err := csvreader.Read()
	if err != nil {
		return nil, err
	}

	err = p.ValidateHeader(read)
	if err != nil {
		return nil, fmt.Errorf("could not read header: %v", err)
	}

	for i, key := range read {
		keys[key] = i
	}

	var queries []QueryParams
	for {
		records, err := csvreader.Read()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return nil, err
		}

		queries = append(queries, p.RowReader(keys, records))
	}

	if len(queries) == 0 {
		return nil, errors.New("provided csv has no queries")
	}

	return queries, nil
}

// RowReader make struct
func (p *Provider) RowReader(keys map[string]int, row []string) QueryParams {
	return QueryParams{
		Host:  row[keys["hostname"]],
		Start: row[keys["start_time"]],
		End:   row[keys["end_time"]],
	}
}

// ExecuteQuery queryWorker query
func (p *Provider) ExecuteQuery(ctx context.Context, conn DbCon, q []QueryParams) ([]Stat, error) {
	//creates the prepared statements
	_, err := conn.Prepare(ctx, configuration.StatementName, configuration.SqlQuery)
	if err != nil {
		return nil, fmt.Errorf("could not prepare statement: %v", err)
	}

	stats := make([]Stat, 0, len(q))
	for _, query := range q {
		// timing the query
		begin := time.Now()
		rows, err := conn.Query(ctx, configuration.StatementName, query.Host, query.Start, query.End)
		if err != nil {
			return nil, fmt.Errorf("could not query: %v", err)
		}

		//iterate rows query
		for rows.Next() {
		}

		rows.Close()
		err = rows.Err()
		if err != nil {
			return nil, fmt.Errorf("could not close rows: %v", err)
		}

		// create slice
		stats = append(stats, Stat{Time: time.Since(begin)})
	}
	return stats, nil
}

// ProcessQueries make workers and send each host to be executed by a worker
func (r *RunQuery) ProcessQueries(ctx context.Context, queries []QueryParams, workers int) (StatSummary, error) {
	hostQueries := map[string][]QueryParams{}
	for _, q := range queries {
		hostQueries[q.Host] = append(hostQueries[q.Host], q)
	}

	c := make(chan []QueryParams, workers)
	stats := make(chan []Stat, workers)

	var g errgroup.Group
	// make each concurrent worker as a go routine
	for i := 0; i < workers; i++ {
		g.Go(r.queryWorker(ctx, c, stats))
	}

	begin := time.Now()

	// send each host value from csv to be executed by a worker
	go func() {
		for _, q := range hostQueries {
			c <- q
		}

		close(c)
	}()

	// close stats channel after all workers are finished
	go func() {
		_ = g.Wait()
		close(stats)
	}()

	//process stats
	totalStats := make([]Stat, 0, len(queries))
	for s := range stats {
		for _, i := range s {
			totalStats = append(totalStats, i)
		}
	}

	return StatSummary{
		Stats:     totalStats,
		TotalTime: time.Since(begin),
	}, g.Wait()

}

func (r *RunQuery) queryWorker(ctx context.Context, c chan []QueryParams, stats chan []Stat) func() error {
	return func() error {
		connection, err := r.Connection(ctx)
		if err != nil {
			return err
		}
		defer connection.Close(ctx)

		for hostnameQueries := range c {
			stat, err := r.Query(ctx, connection, hostnameQueries)
			if err != nil {
				return err
			}
			stats <- stat
		}
		return nil
	}

}

func (p *Provider) RunPerfTest(ctx context.Context, connectionString string, queries []QueryParams) (StatSummary, error) {
	log := logger.SugaredLogger().WithContextCorrelationId(ctx).With("package", "pg", "action", "RunPerfTest")
	log.Debugf("queryWorker perf test")
	runner := &RunQuery{
		Connection: func(ctx context.Context) (DbCon, error) {
			return pgx.Connect(ctx, connectionString)
		},
		Query: p.ExecuteQuery,
	}

	statSummary, err := runner.ProcessQueries(ctx, queries, int(p.Workers))
	if err != nil {
		return statSummary, fmt.Errorf("unable to queryWorker perf test: %v", err)
		//log.Errorf("Unable to queryWorker perf test: %v", err)
	}
	return statSummary, nil
}
