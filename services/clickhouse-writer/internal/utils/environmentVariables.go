package utils

import (
	"log"
	"log/slog"
	"os"
	"strconv"
	"time"
)

func ValidateEnvVars(logger *slog.Logger, vars ...string) {
	missing := []string{}
	for _, v := range vars {
		value := os.Getenv(v)
		if value == "" {
			missing = append(missing, v)
		}
	}
	if len(missing) > 0 {
		logger.Error("Missing required environment variables", "vars", missing)
		os.Exit(1)
	}
}

func GetFlushTimeoutDuration() time.Duration {
	flushTimeoutStr := os.Getenv("CLICK_HOUSE_WRITER_MICROBATCH_FLUSH_TIMEOUT_SECONDS")
	flushTimeoutInt, err := strconv.Atoi(flushTimeoutStr)
	if err != nil {
		log.Fatal("Invalid flush timeout value: ", err)
	}

	return time.Duration(flushTimeoutInt) * time.Second
}
