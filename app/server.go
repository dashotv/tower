package app

import (
	"fmt"

	"github.com/fsnotify/fsnotify"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

type Server struct {
	Router  *gin.Engine
	Log     *zap.SugaredLogger
	watcher *fsnotify.Watcher
}

func (s *Server) Start() error {
	s.Log.Info("starting tower...")

	if err := s.Cron(); err != nil {
		return err
	}

	s.Routes()
	if cfg.Filesystems.Enabled {
		s.Watcher()
		defer s.watcher.Close()
	}

	s.Log.Info("starting web...")
	if err := s.Router.Run(fmt.Sprintf(":%d", cfg.Port)); err != nil {
		return errors.Wrap(err, "starting router")
	}

	return nil
}
