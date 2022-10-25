package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel/attribute"
	oteltrace "go.opentelemetry.io/otel/trace"
	"net/http"
	"path/filepath"
	"postgres-perf/api/response"
	"postgres-perf/postgresql"
	"postgres-perf/utils"
	"postgres-perf/utils/go-stats/concurrency"
	"postgres-perf/utils/logger"
	"strconv"
)

// SaveFileHandler godoc
// @Summary upload csv file
// @Description upload csv file
// @ID uploader
// @Accept  multipart/form-data
// @Produce json
// @Param   file formData file true  "query test file"
// @Success 200 {object} model.JSONSuccessResult "The request was validated and has been processed successfully (sync)"
// @Failure 404 {object} model.JSONFailureResult "The payload was rejected as invalid"
// @Failure 500 {object} model.JSONFailureResult "An internal error has occurred, most likely due to an uncaught exception"
// @Router /v1/upload [post]
func SaveFileHandler(c *gin.Context) {
	concurrency.GlobalWaitGroup.Add(1)
	defer concurrency.GlobalWaitGroup.Done()

	var (
		sqlProvider *postgresql.Provider
		err         error
	)

	log := logger.SugaredLogger().WithContextCorrelationId(c)

	ctx := c.Request.Context()

	if sqlProvider, err = postgresql.New(ctx); err != nil {
		log.Errorf("Error while initializing Sql Provider: %s", err)
		response.FailureResponse(c, nil, utils.HttpError{Code: 500, Err: err})
	}

	file, err := c.FormFile("file")
	if err != nil {
		log.Errorf("Error at file post submit")
		response.FailureResponse(c, nil, utils.HttpError{Code: 400, Err: err})
		return
	}

	// Retrieve file information
	extension := filepath.Ext(file.Filename)
	newFileName := uuid.New().String() + extension
	log.Debugf(file.Filename)
	// TODO file size check and type

	// The file cannot be received.
	if err != nil {
		log.Errorf("Error while parsing request: %s", err.Error())
		response.FailureResponse(c, nil, utils.HttpError{Code: 400, Err: err})
		return
	}

	//The file is received, so let's save it to /tmp
	if err := c.SaveUploadedFile(file, "/tmp/"+newFileName); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "Unable to save the file",
		})
		return
	}

	// query file params
	params, err := sqlProvider.QueryFileParams(ctx, "/tmp/"+newFileName)
	if err != nil {
		return
	}

	// load csv and validate
	queries, err := sqlProvider.LoadCSV(params)
	if err != nil {
		log.Errorf("Unable to parse queries: %s", err)
		response.FailureResponse(c, nil, utils.HttpError{Code: 500, Err: err})
		return
	}

	// connection string
	connectionString := sqlProvider.ConnectionString(ctx)

	// run perf tests
	statSummary, err := sqlProvider.RunPerfTest(ctx, connectionString, queries)
	if err != nil {
		log.Errorf("Unable to run perf  test: %s", err)
		response.FailureResponse(c, nil, utils.HttpError{Code: 500, Err: err})
		return
	}

	// time values
	tq := statSummary.Number()
	tt := statSummary.TotalTime
	at := statSummary.Aggregate()
	mt := statSummary.Max()
	mm := statSummary.Min()

	log.Debugf("total queris: %d, total time: %s, aggregate time: %s, max time: %s, min time: %s",
		tq, tt, at, mt, mm)

	// map
	putbody := map[string]interface{}{
		"totalqueries": tq,
		"totaltime":    tt.String(),
		"aggregate":    at.String(),
		"maxtime":      mt.String(),
		"mintime":      mm.String(),
	}

	// tracer
	_, span := tracer.Start(c.Request.Context(), "PgPerfResult",
		oteltrace.WithAttributes(attribute.String("totalqueries", strconv.Itoa(tq)),
			attribute.String("totaltime", tt.String()),
			attribute.String("aggregate", at.String()),
			attribute.String("maxtime", mt.String()),
			attribute.String("mintime", mm.String()),
			attribute.String("Correlation_id", c.MustGet("correlation_id").(string))))

	defer span.End()

	response.SuccessResponse(c, c.MustGet("correlation_id").(string), putbody)

}
