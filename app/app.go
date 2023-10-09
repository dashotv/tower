package app

import (
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/clerkinc/clerk-sdk-go/clerk"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"github.com/philippgille/gokv/redis"
	"go.uber.org/zap"
)

var once sync.Once
var instance *Application

type Application struct {
	Config *Config
	Router *gin.Engine
	DB     *Connector
	Cache  *Cache
	Log    *zap.SugaredLogger
	// Add additional clients and connections
}

func logger() (*zap.Logger, error) {
	return zap.NewProduction()
}

func initialize() *Application {
	cfg := ConfigInstance()

	zapcfg := zap.NewProductionConfig()

	switch cfg.Mode {
	case "dev":
		zapcfg.Development = true
		zapcfg.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	case "release":
		gin.SetMode(cfg.Mode)
		zapcfg.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	}

	logger, err := logger()
	if err != nil {
		fmt.Printf("logger: %s", err)
		return nil
	}
	log := logger.Sugar()

	db, err := NewConnector()
	if err != nil {
		log.Errorf("database connection failed: %s", err)
	}

	router := gin.New()
	router.Use(ginzap.Ginzap(logger, time.RFC3339, true), ginzap.RecoveryWithZap(logger, true))

	if cfg.Mode == "release" {
		clerKey := os.Getenv("CLERK_SECRET_KEY")
		if clerKey == "" {
			log.Fatal("CLERK_SECRET_KEY is not set")
		}

		clerkClient, err := clerk.NewClient(clerKey)
		if err != nil {
			log.Fatalf("clerk: %s", err)
		}

		router.Use(requireSession(clerkClient))
	}

	cache, err := NewCache(redis.Options{Address: cfg.Redis.Address})
	if err != nil {
		log.Fatalf("cache: %s", err)
	}

	// Add additional clients and connections

	return &Application{
		Config: cfg,
		Router: router,
		DB:     db,
		Cache:  cache,
		Log:    log,
	}
}

func App() *Application {
	once.Do(func() {
		instance = initialize()
	})
	return instance
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
