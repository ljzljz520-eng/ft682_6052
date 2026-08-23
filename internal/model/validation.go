package model

import "strings"

func ValidateVenue(venue Venue) error {
	if strings.TrimSpace(venue.ID) == "" {
		return ErrInvalid("venue id is required")
	}
	if strings.TrimSpace(venue.Name) == "" {
		return ErrInvalid("venue name is required")
	}
	if venue.Capacity <= 0 {
		return ErrInvalid("venue capacity must be positive")
	}
	return nil
}

func ValidateTimeSlot(slot TimeSlot) error {
	if strings.TrimSpace(slot.ID) == "" || strings.TrimSpace(slot.VenueID) == "" {
		return ErrInvalid("slot identity is required")
	}
	if strings.TrimSpace(slot.StartsAt) == "" || strings.TrimSpace(slot.EndsAt) == "" {
		return ErrInvalid("slot range is required")
	}
	if slot.StartsAt >= slot.EndsAt {
		return ErrInvalid("slot range is inverted")
	}
	return nil
}

func ValidateReservation(reservation Reservation) error {
	if strings.TrimSpace(reservation.ID) == "" {
		return ErrInvalid("reservation id is required")
	}
	if strings.TrimSpace(reservation.VenueID) == "" || strings.TrimSpace(reservation.SlotID) == "" {
		return ErrInvalid("venue and slot are required")
	}
	if strings.TrimSpace(reservation.Applicant) == "" {
		return ErrInvalid("applicant is required")
	}
	if strings.TrimSpace(reservation.Purpose) == "" {
		return ErrInvalid("purpose is required")
	}
	if reservation.Status == "" {
		return ErrInvalid("status is required")
	}
	return nil
}

func ValidateReview(status string, detail string) error {
	if status != "approved" && status != "rejected" {
		return ErrInvalid("review status must be approved or rejected")
	}
	if strings.TrimSpace(detail) == "" {
		return ErrInvalid("review detail is required")
	}
	return nil
}

func ValidateCollaboration(note CollaborationNote) error {
	if strings.TrimSpace(note.ReservationID) == "" || strings.TrimSpace(note.Author) == "" {
		return ErrInvalid("note identity is required")
	}
	if strings.TrimSpace(note.Body) == "" {
		return ErrInvalid("note body is required")
	}
	if note.Visibility == "" {
		note.Visibility = "team"
	}
	return nil
}

func ValidateArchive(reason string, actor string) error {
	if strings.TrimSpace(reason) == "" {
		return ErrInvalid("archive reason is required")
	}
	if strings.TrimSpace(actor) == "" {
		return ErrInvalid("archive actor is required")
	}
	return nil
}

func ValidatePage(request PageRequest) PageRequest {
	if request.Page < 1 {
		request.Page = 1
	}
	if request.PageSize < 1 {
		request.PageSize = 20
	}
	if request.PageSize > 100 {
		request.PageSize = 100
	}
	return request
}

func ValidateScope(scope string) string {
	value := strings.ToLower(strings.TrimSpace(scope))
	if value == "" {
		return "all"
	}
	if value != "team" && value != "public" && value != "all" {
		return "all"
	}
	return value
}

func ErrInvalid(message string) error {
	return &ValidationError{Message: message}
}

type ValidationError struct {
	Message string
}

func (e *ValidationError) Error() string {
	return e.Message
}
