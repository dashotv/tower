package app

import (
	"fmt"
	"net/http"
	"time"

	"github.com/clerkinc/clerk-sdk-go/clerk"
	"github.com/fsnotify/fsnotify"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

var server *Server

func setupServer() (err error) {
	if cfg.Mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	server = &Server{
		Log: log.Named("server"),
	}

	server.Engine = gin.New()
	server.Engine.Use(ginzap.Ginzap(log.Desugar(), time.RFC3339, true), ginzap.RecoveryWithZap(log.Desugar(), true))
	server.Default = server.Engine.Group("/")
	server.Router = server.Engine.Group("/")

	server.Routes()

	if cfg.Auth {
		clerkSecret := cfg.ClerkSecretKey
		if clerkSecret == "" {
			log.Fatal("CLERK_SECRET_KEY is not set")
		}

		clerkClient, err := clerk.NewClient(clerkSecret)
		if err != nil {
			log.Fatalf("clerk: %s", err)
		}

		server.Router.Use(requireSession(clerkClient))
	}

	return nil
}

type Server struct {
	Engine  *gin.Engine
	Router  *gin.RouterGroup
	Default *gin.RouterGroup
	Log     *zap.SugaredLogger
	watcher *fsnotify.Watcher
}

func (s *Server) Start() error {
	s.Log.Info("starting tower...")

	go events.Start()
	go func() {
		s.Log.Infof("starting workers (%d)...", cfg.MinionConcurrency)
		workers.Start()
	}()

	s.Log.Info("starting web...")
	if err := s.Engine.Run(fmt.Sprintf(":%d", cfg.Port)); err != nil {
		return errors.Wrap(err, "starting router")
	}

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
