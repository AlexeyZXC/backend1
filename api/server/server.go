package server

import (
	"context"
	"net/http"
	"time"

	"github.com/AlexeyZXC/backend1/tree/CourseProject/app/repo/link"
)

type Server struct {
	srv http.Server
	ls  *link.Links
}

func NewServer(addr string, h http.Handler) *Server {
	s := &Server{}

	s.srv = http.Server{
		Addr:              addr,
		Handler:           h,
		ReadTimeout:       30 * time.Second,
		WriteTimeout:      30 * time.Second,
		ReadHeaderTimeout: 30 * time.Second,
	}

	return s
}

func (s *Server) Stop() {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	s.srv.Shutdown(ctx)
	cancel()
}

func (s *Server) Start(ls *link.Links) {
	s.ls = ls
	go s.srv.ListenAndServe()
}
