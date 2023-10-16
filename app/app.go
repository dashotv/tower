package app

import (
	"net/http"
	"os"
	"time"

	"github.com/clerkinc/clerk-sdk-go/clerk"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"github.com/philippgille/gokv/redis"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var initialized bool
var cfg *Config
var log *zap.SugaredLogger
var logger *zap.Logger
var db *Connector
var router *gin.Engine
var routerDefault *gin.RouterGroup
var routerAuth *gin.RouterGroup
var cache *Cache
var server *Server
var minion *Minion

type SetupFunc func() error

func Start() error {
	err := setup(
		setupConfig,
		setupLogger,
		setupDb,
		setupRouter,
		setupCache,
		setupWorkers,
		setupEvents,
	)
	if err != nil {
		return err
	}

	initialized = true
	log.Info("initialized: ", initialized)
	log.Debugf("config: %+v", cfg)

	server = &Server{Log: log.Named("server"), Engine: router, Router: routerAuth, Default: routerDefault}
	return server.Start()
}

func setup(fs ...SetupFunc) error {
	for _, f := range fs {
		if err := f(); err != nil {
			return err
		}
	}

	return nil
}

func setupConfig() (err error) {
	cfg = &Config{}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err != nil {
		// fmt.Printf("WARN: unable to read config: %s\n", err)
		return nil //errors.Wrap(err, "unable to read config")
	}

	if err := viper.Unmarshal(cfg); err != nil {
		return errors.Wrap(err, "failed to unmarshal configuration file")
	}

	if err := cfg.Validate(); err != nil {
		return errors.Wrap(err, "failed to validate config")
	}

	return nil
}

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

func setupCache() (err error) {
	cache, err = NewCache(log.Named("cache"), redis.Options{Address: cfg.Redis.Address})
	if err != nil {
		return err
	}

	return nil
}

func setupDb() (err error) {
	db, err = NewConnector()
	if err != nil {
		return err
	}

	return nil
}

func setupRouter() (err error) {
	if cfg.Mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	router = gin.New()
	router.Use(ginzap.Ginzap(logger, time.RFC3339, true), ginzap.RecoveryWithZap(logger, true))
	routerDefault = router.Group("/")
	routerAuth = router.Group("/")

	if cfg.Auth {
		clerKey := os.Getenv("CLERK_SECRET_KEY")
		if clerKey == "" {
			log.Fatal("CLERK_SECRET_KEY is not set")
		}

		clerkClient, err := clerk.NewClient(clerKey)
		if err != nil {
			log.Fatalf("clerk: %s", err)
		}

		routerAuth.Use(requireSession(clerkClient))
	}

	return nil
}

func setupWorkers() error {
	minion = NewMinion(cfg.Minion.Concurrency)
	return nil
}

// AsGin converts middleware to the gin middleware handler.
func requireSession(client clerk.Client) gin.HandlerFunc {
	requireActionSession := clerk.RequireSessionV2(client)
	return func(gctx *gin.Context) {
		var skip = true
		var handler http.HandlerFunc = func(http.ResponseWriter, *http.Request) {
			skip = false
		}
		requireActionSession(handler).ServeHTTP(gctx.Writer, gctx.Request)
		switch {
		case skip:
			gctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "session required"})
		default:
			gctx.Next()
		}
	}
}
