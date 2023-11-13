package http

import (
	"context"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"net"
	"net/http"
	"time"
)

type routerHandler interface {
	GetHandler(eh *EndpointHandler) http.Handler
}

type Server struct {
	logger          *zap.SugaredLogger
	shutdownTimeout time.Duration
	client          *http.Server
	listener        net.Listener
	isReady         bool
	EndpointHandler *EndpointHandler
}

func NewServer(port int, shutdownTimeout time.Duration,
	routerHandler routerHandler, logger *zap.SugaredLogger,
	endpoinHandler *EndpointHandler) (*Server, error) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, fmt.Errorf("cannot bind HTTP server '%d': %v", port, err)
	}
	return &Server{
		client: &http.Server{
			Handler: routerHandler.GetHandler(endpoinHandler),
		},
		listener:        listener,
		logger:          logger,
		shutdownTimeout: shutdownTimeout,
		isReady:         false,
	}, nil
}

func (s *Server) Ready() error {
	if s.isReady {
		return nil
	}

	return errors.New("I am not ready")
}

func (s *Server) Stop() error {
	s.isReady = false
	s.logger.Infof("[%s] HTTP server is stoping...", s.listener.Addr().String())

	ctx, cancel := context.WithTimeout(context.Background(), s.shutdownTimeout)
	defer cancel()

	s.client.SetKeepAlivesEnabled(false)

	if err := s.client.Shutdown(ctx); err != nil {
		return fmt.Errorf("cannot stop HTTP server: %w", err)
	}
	s.logger.Infof("[%s] HTTP server was stopped", s.listener.Addr().String())

	return nil
}

func (s *Server) Run() {
	s.logger.Infof("[%s] HTTP server is running...", s.listener.Addr().String())

	go func() {
		s.isReady = true
		s.logger.Infof("[%s] HTTP server is run", s.listener.Addr().String())

		if err := s.client.Serve(s.listener); err != nil {
			s.isReady = false
			if errors.Is(err, http.ErrServerClosed) {
				return
			}

			s.logger.Errorf("[%s] HTTP server was stopped with error: %s", s.listener.Addr().String(), err)
		}
	}()
}
