package server

import (
	"context"
	"log"

	"github.com/labstack/echo"
)

type Config struct {
	Port string
}

type Server struct {
	*echo.Echo
	cfg Config
}

func New(cfg Config) *Server {
	return &Server{
		cfg:  cfg,
		Echo: echo.New(),
	}
}

func (s *Server) Start(context.Context) error {
	log.Println(s.Echo.Start(":" + s.cfg.Port))

	return nil
}

func (s *Server) ShutDown(ctx context.Context) error {
	err := s.Echo.Shutdown(ctx)
	if err != nil {
		return err
	}

	return nil
}
