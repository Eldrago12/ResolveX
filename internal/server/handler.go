package server

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func (s *Server) resolveHandler(w http.ResponseWriter, r *http.Request) {
	domain := r.URL.Query().Get("domain")
	if domain == "" {
		http.Error(w, "domain parameter is required", http.StatusBadRequest)
		return
	}

	clientIP := r.RemoteAddr
	allowed, err := s.Limiter.Allow(r.Context(), clientIP)
	if err != nil || !allowed {
		http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
		return
	}

	ip, err := s.Resolver.ResolveDomain(r.Context(), domain)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error resolving domain: %v", err), http.StatusInternalServerError)
		return
	}

	response := map[string]string{"domain": domain, "ip": ip}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
