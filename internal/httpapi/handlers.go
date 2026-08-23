package httpapi

import (
	"net/http"
	"strconv"
	"strings"

	"venue-reservation/internal/model"
)

func (s *Server) handleRoot(writer http.ResponseWriter, request *http.Request) {
	if request.URL.Path != "/" {
		writeJSON(writer, http.StatusNotFound, map[string]string{"error": "route not found"})
		return
	}
	writeJSON(writer, http.StatusOK, map[string]any{"service": "venue-reservation", "status": "ready"})
}

func (s *Server) handleHealth(writer http.ResponseWriter, request *http.Request) {
	if !jsonMethod(request, http.MethodGet) {
		writeJSON(writer, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
		return
	}
	writeJSON(writer, http.StatusOK, map[string]string{"status": "ok"})
}

func (s *Server) handleVenues(writer http.ResponseWriter, request *http.Request) {
	if request.Method == http.MethodGet {
		onlyEnabled := request.URL.Query().Get("enabled") == "true"
		venues, err := s.service.ListVenues(onlyEnabled)
		if err != nil {
			writeError(writer, err)
			return
		}
		writeJSON(writer, http.StatusOK, map[string]any{"items": venues})
		return
	}
	if request.Method != http.MethodPost {
		writeJSON(writer, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
		return
	}
	var input venueRequest
	if err := decodeJSON(request, &input); err != nil {
		writeError(writer, err)
		return
	}
	if err := s.service.RegisterVenue(venueFromRequest(input)); err != nil {
		writeError(writer, err)
		return
	}
	writeJSON(writer, http.StatusCreated, venueFromRequest(input))
}

func (s *Server) handleSlots(writer http.ResponseWriter, request *http.Request) {
	if request.Method == http.MethodGet {
		slots, err := s.service.ListTimeSlots(request.URL.Query().Get("venue_id"), request.URL.Query().Get("available") == "true")
		if err != nil {
			writeError(writer, err)
			return
		}
		writeJSON(writer, http.StatusOK, map[string]any{"items": slots})
		return
	}
	if request.Method != http.MethodPost {
		writeJSON(writer, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
		return
	}
	var input slotRequest
	if err := decodeJSON(request, &input); err != nil {
		writeError(writer, err)
		return
	}
	slot := slotFromRequest(input)
	if err := s.service.AddTimeSlot(slot); err != nil {
		writeError(writer, err)
		return
	}
	writeJSON(writer, http.StatusCreated, slot)
}

func (s *Server) handleReservations(writer http.ResponseWriter, request *http.Request) {
	if request.Method == http.MethodGet {
		page := parseInt(request.URL.Query().Get("page"), 1)
		pageSize := parseInt(request.URL.Query().Get("page_size"), 20)
		var result model.Page[model.Reservation]
		var err error
		member := strings.TrimSpace(request.URL.Query().Get("member"))
		if member != "" {
			result, err = s.service.ListForMember(member, model.PageRequest{Page: page, PageSize: pageSize})
		} else {
			filter := model.QueryFilter(request.URL.Query())
			result, err = s.service.FindReservations(filter, model.PageRequest{Page: page, PageSize: pageSize})
		}
		if err != nil {
			writeError(writer, err)
			return
		}
		writeJSON(writer, http.StatusOK, result)
		return
	}
	if request.Method != http.MethodPost {
		writeJSON(writer, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
		return
	}
	var input reservationRequest
	if err := decodeJSON(request, &input); err != nil {
		writeError(writer, err)
		return
	}
	reservation, err := s.service.SubmitReservation(reservationFromRequest(input))
	if err != nil {
		writeError(writer, err)
		return
	}
	writeJSON(writer, http.StatusCreated, reservation)
}

func (s *Server) handleReservationPath(writer http.ResponseWriter, request *http.Request) {
	parts := pathParts(strings.TrimPrefix(request.URL.Path, "/"))
	if len(parts) < 2 || parts[0] != "reservations" {
		writeJSON(writer, http.StatusNotFound, map[string]string{"error": "reservation route not found"})
		return
	}
	id := parts[1]
	if len(parts) == 2 && request.Method == http.MethodGet {
		detail, err := s.service.GetReservationDetail(id)
		if err != nil {
			writeError(writer, err)
			return
		}
		writeJSON(writer, http.StatusOK, detail)
		return
	}
	if request.Method != http.MethodPost || len(parts) != 3 {
		writeJSON(writer, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
		return
	}
	s.handleReservationAction(writer, request, id, parts[2])
}

func (s *Server) handleReservationAction(writer http.ResponseWriter, request *http.Request, id string, action string) {
	switch action {
	case "review":
		var input reviewRequest
		if err := decodeJSON(request, &input); err != nil {
			writeError(writer, err)
			return
		}
		value, err := s.service.ReviewReservation(id, input.Status, input.Actor, input.Detail)
		if err != nil {
			writeError(writer, err)
			return
		}
		writeJSON(writer, http.StatusOK, value)
	case "notes":
		var input noteRequest
		if err := decodeJSON(request, &input); err != nil {
			writeError(writer, err)
			return
		}
		value, err := s.service.AddCollaboration(id, input.Author, input.Body, input.Visibility)
		if err != nil {
			writeError(writer, err)
			return
		}
		writeJSON(writer, http.StatusCreated, value)
	case "archive":
		var input archiveRequest
		if err := decodeJSON(request, &input); err != nil {
			writeError(writer, err)
			return
		}
		value, err := s.service.ArchiveReservation(id, input.Actor, input.Reason)
		if err != nil {
			writeError(writer, err)
			return
		}
		writeJSON(writer, http.StatusOK, value)
	case "update":
		var input updateRequest
		if err := decodeJSON(request, &input); err != nil {
			writeError(writer, err)
			return
		}
		value, err := s.service.UpdateReservation(id, input.Purpose, input.Notes)
		if err != nil {
			writeError(writer, err)
			return
		}
		writeJSON(writer, http.StatusOK, value)
	default:
		writeJSON(writer, http.StatusNotFound, map[string]string{"error": "reservation action not found"})
	}
}

func (s *Server) handleReportPath(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodGet {
		writeJSON(writer, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
		return
	}
	parts := pathParts(strings.TrimPrefix(request.URL.Path, "/"))
	if len(parts) != 2 || parts[0] != "reports" {
		writeJSON(writer, http.StatusNotFound, map[string]string{"error": "report route not found"})
		return
	}
	report, err := s.service.PublishReport(parts[1])
	if err != nil {
		writeError(writer, err)
		return
	}
	writeJSON(writer, http.StatusOK, report)
}

func (s *Server) handleReservationExport(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodGet {
		writeJSON(writer, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
		return
	}
	data, err := s.service.ExportCSV(model.QueryFilter(request.URL.Query()))
	if err != nil {
		writeError(writer, err)
		return
	}
	writer.Header().Set("Content-Type", "text/csv")
	writer.Header().Set("Content-Disposition", "attachment; filename=reservations.csv")
	writer.WriteHeader(http.StatusOK)
	_, _ = writer.Write(data)
}

func parseInt(value string, fallback int) int {
	parsed, err := strconv.Atoi(value)
	if err != nil || parsed < 1 {
		return fallback
	}
	return parsed
}
