package httpapi

import "net/http"

type systemPage struct {
	Name         string   `json:"name"`
	Version      string   `json:"version"`
	Storage      string   `json:"storage"`
	Capabilities []string `json:"capabilities"`
	Endpoints    []string `json:"endpoints"`
	Guarantees   []string `json:"guarantees"`
}

func (s *Server) handleSystem(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodGet {
		writeJSON(writer, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
		return
	}
	writeJSON(writer, http.StatusOK, systemPage{Name: "Venue Reservation Engineering Skeleton", Version: "1.0", Storage: "go.etcd.io/bbolt", Capabilities: systemCapabilities(), Endpoints: systemEndpoints(), Guarantees: systemGuarantees()})
}

func systemCapabilities() []string {
	return []string{
		"venue registration and availability",
		"deterministic reservation submission",
		"review and approval history",
		"member-scoped search",
		"team and public collaboration",
		"archive and venue reports",
		"validated reservation import",
		"CSV reservation export",
	}
}

func systemEndpoints() []string {
	return []string{"GET /health", "GET /venues", "GET /slots", "GET /reservations", "GET /reservations/{id}", "POST /reservations/{id}/review", "POST /reservations/{id}/notes", "POST /reservations/{id}/archive", "GET /reports/{venue_id}", "GET /exports/reservations", "GET /system"}
}

func systemGuarantees() []string {
	return []string{"local embedded persistence", "no external service dependency", "deterministic seed catalog", "reopen-safe entity reads"}
}
