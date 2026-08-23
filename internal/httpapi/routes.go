package httpapi

import "net/http"

func (s *Server) registerRoutes() {
	s.mux.HandleFunc("/", s.handleRoot)
	s.mux.HandleFunc("/health", s.handleHealth)
	s.mux.HandleFunc("/venues", s.handleVenues)
	s.mux.HandleFunc("/slots", s.handleSlots)
	s.mux.HandleFunc("/reservations", s.handleReservations)
	s.mux.HandleFunc("/reservations/", s.handleReservationPath)
	s.mux.HandleFunc("/reports/", s.handleReportPath)
	s.mux.HandleFunc("/exports/reservations", s.handleReservationExport)
	s.mux.HandleFunc("/system", s.handleSystem)
}

func pathParts(path string) []string {
	parts := make([]string, 0)
	for _, item := range splitPath(path) {
		if item != "" {
			parts = append(parts, item)
		}
	}
	return parts
}

func splitPath(path string) []string {
	result := make([]string, 0)
	start := 0
	for index := 0; index <= len(path); index++ {
		if index == len(path) || path[index] == '/' {
			result = append(result, path[start:index])
			start = index + 1
		}
	}
	return result
}

func jsonMethod(request *http.Request, expected string) bool {
	return request.Method == expected
}
