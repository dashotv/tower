package app

import (
	"os"

	"github.com/streamingfast/logging"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"golang.org/x/term"
)

var log *zap.SugaredLogger
var logger *zap.Logger

func setupLogger() (err error) {
	zapcfg := zap.NewProductionConfig()
	verbosity := 1

	switch cfg.Logger {
	case "dev":
		isTTY := term.IsTerminal(int(os.Stderr.Fd()))
		logStdoutWriter := zapcore.Lock(os.Stderr)
		logger = zap.New(zapcore.NewCore(logging.NewEncoder(verbosity, isTTY), logStdoutWriter, zapcore.DebugLevel))
	case "release":
		zapcfg.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
		logger, err = zap.NewProduction()
		if err != nil {
			return err
		}
	}

	log = logger.Sugar().Named("app")

	return nil
}
