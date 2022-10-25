package postgresql

import "time"

type Stat struct {
	Time time.Duration
}

type StatSummary struct {
	Stats     []Stat
	TotalTime time.Duration
}

// Number total number of queries
func (s StatSummary) Number() int {
	return len(s.Stats)
}

// Aggregate sum time of queries
func (s StatSummary) Aggregate() time.Duration {
	sum := time.Duration(0)
	for _, s := range s.Stats {
		sum = sum + s.Time
	}
	return sum
}

// Max maximum duration of query
func (s StatSummary) Max() time.Duration {
	var max = s.Stats[0].Time
	for _, s := range s.Stats {
		if s.Time > max {
			max = s.Time
		}
	}
	return max
}

// Min minimum duration of query
func (s StatSummary) Min() time.Duration {
	var min = s.Stats[0].Time
	for _, s := range s.Stats {
		if s.Time < min {
			min = s.Time
		}
	}
	return min
}
