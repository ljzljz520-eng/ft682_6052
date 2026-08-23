package httpapi

import (
	"log"
	"net/http"

	"venue-reservation/internal/service"
)

type Server struct {
	service *service.Service
	mux     *http.ServeMux
}

func NewServer(serviceRef *service.Service) *Server {
	server := &Server{service: serviceRef, mux: http.NewServeMux()}
	server.registerRoutes()
	return server
}

func (s *Server) Handler() http.Handler {
	return s.withHeaders(s.mux)
}

func (s *Server) ListenAndServe(address string) error {
	return http.ListenAndServe(address, s.Handler())
}

func (s *Server) withHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-Type", "application/json")
		writer.Header().Set("Cache-Control", "no-store")
		next.ServeHTTP(writer, request)
	})
}

func requestLog(message string) {
	log.Printf("venue-reservation %s", message)
}

func methodAllowed(request *http.Request, methods ...string) bool {
	for _, method := range methods {
		if request.Method == method {
			return true
		}
	}
	return false
}
