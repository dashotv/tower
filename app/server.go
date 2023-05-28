package app

import (
	"fmt"

	"github.com/fsnotify/fsnotify"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type Server struct {
	Router  *gin.Engine
	Log     *logrus.Entry
	watcher *fsnotify.Watcher
}

func New() (*Server, error) {
	log := App().Log.WithField("prefix", "server")
	s := &Server{Log: log, Router: App().Router}

	return s, nil
}

func (s *Server) Start() error {
	s.Log.Info("starting tower...")

	if err := s.Cron(); err != nil {
		return err
	}

	s.Routes()
	if ConfigInstance().Filesystems.Enabled {
		s.Watcher()
		defer s.watcher.Close()
	}

	s.Log.Info("starting web...")
	if err := s.Router.Run(fmt.Sprintf(":%d", cfg.Port)); err != nil {
		return errors.Wrap(err, "starting router")
	}

	return nil
}
