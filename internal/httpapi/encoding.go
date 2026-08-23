package httpapi

import (
	"encoding/json"
	"errors"
	"net/http"

	"venue-reservation/internal/model"
	"venue-reservation/internal/service"
	"venue-reservation/internal/store"
)

type venueRequest struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Address  string `json:"address"`
	Capacity int    `json:"capacity"`
	Enabled  bool   `json:"enabled"`
}

type slotRequest struct {
	ID       string `json:"id"`
	VenueID  string `json:"venue_id"`
	StartsAt string `json:"starts_at"`
	EndsAt   string `json:"ends_at"`
}

type reservationRequest struct {
	ID        string `json:"id"`
	VenueID   string `json:"venue_id"`
	SlotID    string `json:"slot_id"`
	Applicant string `json:"applicant"`
	Purpose   string `json:"purpose"`
	Scope     string `json:"scope"`
	Notes     string `json:"notes"`
}

type reviewRequest struct {
	Status string `json:"status"`
	Actor  string `json:"actor"`
	Detail string `json:"detail"`
}

type noteRequest struct {
	Author     string `json:"author"`
	Body       string `json:"body"`
	Visibility string `json:"visibility"`
}

type archiveRequest struct {
	Actor  string `json:"actor"`
	Reason string `json:"reason"`
}

type updateRequest struct {
	Purpose string `json:"purpose"`
	Notes   string `json:"notes"`
}

func decodeJSON(request *http.Request, target any) error {
	decoder := json.NewDecoder(request.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(target); err != nil {
		return model.ErrInvalid("request body is invalid")
	}
	return nil
}

func writeJSON(writer http.ResponseWriter, status int, value any) {
	writer.WriteHeader(status)
	if err := json.NewEncoder(writer).Encode(value); err != nil {
		requestLog("response encoding failed")
	}
}

func writeError(writer http.ResponseWriter, err error) {
	status := http.StatusInternalServerError
	if errors.Is(err, store.ErrNotFound) {
		status = http.StatusNotFound
	} else {
		var validation *model.ValidationError
		if errors.As(err, &validation) {
			status = http.StatusBadRequest
		}
		if service.IsNotFound(err) {
			status = http.StatusNotFound
		}
	}
	writeJSON(writer, status, map[string]string{"error": err.Error()})
}

func venueFromRequest(input venueRequest) model.Venue {
	return model.Venue{ID: input.ID, Name: input.Name, Address: input.Address, Capacity: input.Capacity, Enabled: input.Enabled}
}

func slotFromRequest(input slotRequest) model.TimeSlot {
	return model.TimeSlot{ID: input.ID, VenueID: input.VenueID, StartsAt: input.StartsAt, EndsAt: input.EndsAt}
}

func reservationFromRequest(input reservationRequest) model.Reservation {
	return model.Reservation{ID: input.ID, VenueID: input.VenueID, SlotID: input.SlotID, Applicant: input.Applicant, Purpose: input.Purpose, Scope: input.Scope, Notes: input.Notes}
}
