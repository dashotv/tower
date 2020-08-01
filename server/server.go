package server

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/dashotv/tower/application"
	"github.com/dashotv/tower/config"
)

var cfg *config.Config
var app *application.App

type Server struct {
	Router *gin.Engine
	Log    *logrus.Entry
}

func New() (*Server, error) {
	cfg = config.Instance()
	app = application.Instance()
	log := app.Log.WithField("prefix", "server")
	s := &Server{Log: log, Router: app.Router}

	return s, nil
}

func (s *Server) Start() error {
	s.Log.Info("starting tower...")

	s.Routes()

	//s.Jobs configuration

	s.Log.Info("starting web...")
	if err := s.Router.Run(fmt.Sprintf(":%d", cfg.Port)); err != nil {
		return errors.Wrap(err, "starting router")
	}

	return nil
}
