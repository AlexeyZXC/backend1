package server

import (
	"context"
	"net/http"
	"sync"
	"time"

	"github.com/AlexeyZXC/backend1/tree/CourseProject/app/repo/link"
	"go.uber.org/zap"
)

type Server struct {
	srv http.Server
	ls  *link.Links
	wg  *sync.WaitGroup
	log *zap.SugaredLogger
}

func NewServer(addr string, h http.Handler, wg *sync.WaitGroup, log *zap.SugaredLogger) *Server {
	s := &Server{}

	s.srv = http.Server{
		Addr:              addr,
		Handler:           h,
		ReadTimeout:       30 * time.Second,
		WriteTimeout:      30 * time.Second,
		ReadHeaderTimeout: 30 * time.Second,
	}
	s.wg = wg
	s.log = log

	return s
}

func (s *Server) Stop() {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	s.srv.Shutdown(ctx)
	cancel()
}

func (s *Server) Start(ls *link.Links) {
	s.ls = ls
	go func() {
		s.log.Infof("--- Working ---")
		err := s.srv.ListenAndServe()
		if err != nil {
			s.log.Infof("Failed to start server: ", err)
		}
		s.wg.Add(-1)
	}()
}
