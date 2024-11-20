package server

import (
	"context"
	"log"
	"net/http"

	"github.com/Eldrago12/ResolveX/internal/prefetch"

	"github.com/Eldrago12/ResolveX/internal/resolver"

	"github.com/Eldrago12/ResolveX/internal/limiter"
)

type Server struct {
	Resolver   *resolver.Resolver
	Limiter    *limiter.RateLimiter
	Prefetcher *prefetch.Prefetcher
	Port       string
}

func (s *Server) Run() {
	http.HandleFunc("/resolve", s.resolveHandler)
	log.Printf("Server running on port %s\n", s.Port)
	log.Fatal(http.ListenAndServe(":"+s.Port, nil))
}

func (s *Server) StartPrefetcher(ctx context.Context) {
	go s.Prefetcher.Start(ctx)
}
