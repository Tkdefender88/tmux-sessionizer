package logging

import (
	"io"
	"log/slog"
	"os"
	"time"

	"github.com/charmbracelet/log"
	"github.com/spf13/viper"
)

func SetupLogging(w io.Writer) {

	cfg := viper.GetViper()
	isDebug := cfg.GetBool("debug")

	var logWriter io.Writer
	if isDebug {
		logWriter = io.MultiWriter(w, os.Stdout)
	} else {
		logWriter = w
	}

	logHandler := log.NewWithOptions(logWriter, log.Options{
		TimeFormat: time.Kitchen,
		Level:      log.DebugLevel,
	})

	logger := slog.New(logHandler)
	slog.SetDefault(logger)
}
