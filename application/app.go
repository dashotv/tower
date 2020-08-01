package application

import (
	"os"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	ginlogrus "github.com/toorop/gin-logrus"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"

	"github.com/dashotv/tower/config"
	"github.com/dashotv/tower/models"
)

var once sync.Once
var instance *App

type App struct {
	Config *config.Config
	Router *gin.Engine
	// Cache  *redis.Client
	Log *logrus.Entry
	// Add additional clients and connections
	DB *models.Connector
}

func logger() *logrus.Entry {
	logrus.SetLevel(logrus.InfoLevel)
	logrus.SetFormatter(&prefixed.TextFormatter{DisableTimestamp: false, FullTimestamp: true})
	host, _ := os.Hostname()
	return logrus.WithField("prefix", host)
}

func initialize() *App {
	cfg := config.Instance()
	log := logger()

	router := gin.New()
	router.Use(ginlogrus.Logger(log), gin.Recovery())

	db, err := models.NewConnector()
	if err != nil {
		log.Errorf("database connection failed: %s", err)
	}

	if cfg.Mode == "dev" {
		logrus.SetLevel(logrus.DebugLevel)
	}

	if cfg.Mode == "release" {
		gin.SetMode(cfg.Mode)
	}

	// TODO: add this to config
	// cache := redis.NewClient(&redis.Options{
	//	Addr: "localhost:6379",
	//	DB:   15, // use default DB
	// })

	// Add additional clients and connections

	return &App{
		Config: cfg,
		Router: router,
		// Cache:    cache,
		Log: log,
		DB:  db,
	}
}

func Instance() *App {
	once.Do(func() {
		instance = initialize()
	})
	return instance
}
