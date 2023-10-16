package app

import "go.uber.org/zap"

var log *zap.SugaredLogger
var logger *zap.Logger

func setupLogger() (err error) {
	zapcfg := zap.NewProductionConfig()

	switch cfg.Mode {
	case "dev":
		zapcfg.Development = true
		zapcfg.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	case "release":
		zapcfg.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	}

	logger, err = zap.NewProduction()
	if err != nil {
		return err
	}
	log = logger.Sugar().Named("app")

	return nil
}
