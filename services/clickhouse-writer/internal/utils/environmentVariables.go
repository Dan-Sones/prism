package utils

import (
	"log/slog"
	"os"
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
