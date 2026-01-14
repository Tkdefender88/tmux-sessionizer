package logging

import (
	"io"
	"log/slog"
	"os"
	"time"

	"github.com/charmbracelet/log"
	"github.com/muesli/termenv"
	"github.com/spf13/viper"
)

func SetupLogging(w io.Writer) *slog.Logger {

	cfg := viper.GetViper()
	isDebug := cfg.GetBool("debug")

	var logWriter io.Writer
	if isDebug {
		logWriter = io.MultiWriter(w, os.Stdout)
	} else {
		logWriter = w
	}

	logHandler := log.NewWithOptions(logWriter, log.Options{
		Formatter:       log.TextFormatter,
		TimeFormat:      time.Kitchen,
		ReportTimestamp: true,
		ReportCaller:    true,
		Level:           log.DebugLevel,
	})

	// Force color output when in debug mode (writing to terminal)
	if isDebug {
		logHandler.SetColorProfile(termenv.TrueColor)
	}

	logger := slog.New(logHandler)
	slog.SetDefault(logger)

	return logger
}
