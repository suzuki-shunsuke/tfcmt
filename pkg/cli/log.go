package cli

import (
	"log/slog"

	"github.com/suzuki-shunsuke/slog-util/slogutil"
)

func setLogLevel(logLevelVar *slog.LevelVar, logLevel string) {
	if logLevel == "" {
		return
	}
	if err := slogutil.SetLevel(logLevelVar, logLevel); err != nil {
		slog.Error("the log level is invalid", "log_level", logLevel, "error", err)
	}
}
