package app

import (
	"os"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	ginlogrus "github.com/toorop/gin-logrus"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
)

var once sync.Once
var instance *Application

type Application struct {
	Config *Config
	Router *gin.Engine
	DB     *Connector
	// Cache  *redis.Client
	Log *logrus.Entry
	// Add additional clients and connections
}

func logger() *logrus.Entry {
	logrus.SetLevel(logrus.InfoLevel)
	logrus.SetFormatter(&prefixed.TextFormatter{DisableTimestamp: false, FullTimestamp: true})
	host, _ := os.Hostname()
	return logrus.WithField("prefix", host)
}

func initialize() *Application {
	cfg := ConfigInstance()
	log := logger()

	db, err := NewConnector()
	if err != nil {
		log.Errorf("database connection failed: %s", err)
	}

	if cfg.Mode == "dev" {
		logrus.SetLevel(logrus.DebugLevel)
	}

	if cfg.Mode == "release" {
		gin.SetMode(cfg.Mode)
	}

	router := gin.New()
	router.Use(ginlogrus.Logger(log), gin.Recovery())

	// TODO: add this to config
	// cache := redis.NewClient(&redis.Options{
	//	Addr: "localhost:6379",
	//	DB:   15, // use default DB
	// })

	// Add additional clients and connections

	return &Application{
		Config: cfg,
		Router: router,
		DB:     db,
		// Cache:    cache,
		Log: log,
	}
}

func App() *Application {
	once.Do(func() {
		instance = initialize()
	})
	return instance
}
