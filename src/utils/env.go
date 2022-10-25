package utils

import (
	"os"
	"strconv"
	"strings"
)

// EnvOrDefault returns the value of the environment variable identified by name
// parameter, or the value of the def parameter if the environment variable is
// missing.
//
// Note: It will return an empty string if the environment variable is set to an
// empty string.
func EnvOrDefault(name, def string) string {
	if v, ok := os.LookupEnv(name); ok && v != "" {
		return v
	}
	return def
}

// EnvOrDefaultInt64 returns the value of the environment variable identified by
// name parameter cast to a 64bit decimal integer, or the value of the def
// parameter if the environment variable is missing or empty.
//
// Note: It will return the value of def if the environment variable is set to an empty string.
func EnvOrDefaultInt64(name string, def int64) int64 {
	if v, ok := os.LookupEnv(name); ok && v != "" {
		vc, _ := strconv.ParseInt(v, 10, 32)
		return vc
	}
	return def
}

// EnvOrDefaultInt32 returns the value of the environment variable identified by
// name parameter cast to a 32bit decimal integer, or the value of the def
// parameter if the environment variable is missing or empty.
//
// Note: It will return the value of def if the environment variable is set to an empty string.
func EnvOrDefaultInt32(name string, def int32) int32 {
	if v, ok := os.LookupEnv(name); ok && v != "" {
		vc, _ := strconv.ParseInt(v, 10, 32)
		return int32(vc)
	}
	return def
}

// EnvOrDefaultBool returns the value of the environment variable identified by
// name parameter cast to a boolean, or the value of the def parameter if the
// environment variable is missing or empty.
//
// Note: It will return the value of def if the environment variable is set to an empty string.
func EnvOrDefaultBool(name string, def bool) bool {
	if v, ok := os.LookupEnv(name); ok && v != "" {
		vc, _ := strconv.ParseBool(v)
		return vc
	}
	return def
}

// EnvOrDefaultStringSlice returns the value of the environment variable
// identified by name parameter split by the value of the separator parameter as
// a slice of strings, or the value of the def parameter if the environment
// variable is missing.
//
// Note: It will return an empty slice if the environment variable is set to an empty string.
func EnvOrDefaultStringSlice(name, separator string, def []string) []string {
	if v, ok := os.LookupEnv(name); ok && v != "" {
		var result []string
		for _, substr := range strings.Split(v, separator) {
			trimmed := strings.TrimSpace(substr)
			if trimmed == "" {
				continue
			}
			result = append(result, trimmed)
		}
		return result
	}
	return def
}
