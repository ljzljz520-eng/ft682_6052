package service

import (
	"errors"
	"strings"

	"venue-reservation/internal/access"
	"venue-reservation/internal/model"
	"venue-reservation/internal/store"
	"venue-reservation/internal/timeline"
)

func (s *Service) SubmitReservation(reservation model.Reservation) (model.Reservation, error) {
	reservation.Applicant = strings.TrimSpace(reservation.Applicant)
	reservation.Purpose = strings.TrimSpace(reservation.Purpose)
	reservation.Scope = model.ValidateScope(reservation.Scope)
	if _, err := s.GetVenue(reservation.VenueID); err != nil {
		return model.Reservation{}, err
	}
	slot, err := s.GetTimeSlot(reservation.SlotID)
	if err != nil {
		return model.Reservation{}, err
	}
	if !slot.Available || slot.State != "open" {
		return model.Reservation{}, errors.New("slot is unavailable")
	}
	if reservation.ID == "" {
		reservation.ID, _, err = s.store.AllocateID("RSV")
		if err != nil {
			return model.Reservation{}, err
		}
	}
	reservation.Status = "pending"
	reservation.CreatedAt = "sequence"
	reservation.UpdatedAt = reservation.CreatedAt
	reservation.Revision = 1
	if err := model.ValidateReservation(reservation); err != nil {
		return model.Reservation{}, err
	}
	if err := s.store.PutReservation(reservation); err != nil {
		return model.Reservation{}, err
	}
	if err := s.appendEvent(reservation.ID, "submitted", reservation.Applicant, "reservation submitted"); err != nil {
		return model.Reservation{}, err
	}
	return reservation, nil
}

func (s *Service) FindReservations(filter model.ReservationFilter, request model.PageRequest) (model.Page[model.Reservation], error) {
	filter.Scope = model.ValidateScope(filter.Scope)
	filter.Query = strings.TrimSpace(filter.Query)
	return s.store.ListReservations(filter, request)
}

func (s *Service) GetReservation(id string) (model.Reservation, error) {
	return s.store.GetReservation(strings.TrimSpace(id))
}

func (s *Service) GetReservationDetail(id string) (model.ReservationDetail, error) {
	reservation, err := s.GetReservation(id)
	if err != nil {
		return model.ReservationDetail{}, err
	}
	venue, err := s.GetVenue(reservation.VenueID)
	if err != nil {
		return model.ReservationDetail{}, err
	}
	slot, err := s.GetTimeSlot(reservation.SlotID)
	if err != nil {
		return model.ReservationDetail{}, err
	}
	events, err := s.store.ListEvents(reservation.ID)
	if err != nil {
		return model.ReservationDetail{}, err
	}
	notes, err := s.store.ListNotes(reservation.ID)
	if err != nil {
		return model.ReservationDetail{}, err
	}
	archive, err := s.store.GetArchive(reservation.ID)
	if err != nil {
		return model.ReservationDetail{}, err
	}
	activity := timeline.Build(reservation, events, notes, archive)
	view := make([]model.TimelineEntry, 0, len(activity))
	for _, entry := range activity {
		view = append(view, model.TimelineEntry{ID: entry.ID, Kind: entry.Kind, Actor: entry.Actor, Summary: entry.Summary, Sequence: entry.Sequence, Visibility: entry.Visibility})
	}
	return model.ReservationDetail{Reservation: reservation, Venue: venue, Slot: slot, Events: events, Notes: notes, Archive: archive, Timeline: view}, nil
}

func (s *Service) UpdateReservation(id string, purpose string, notes string) (model.Reservation, error) {
	reservation, err := s.GetReservation(id)
	if err != nil {
		return model.Reservation{}, err
	}
	if model.StatusIsTerminal(reservation.Status) {
		return model.Reservation{}, errors.New("terminal reservation cannot be edited")
	}
	if strings.TrimSpace(purpose) == "" {
		return model.Reservation{}, model.ErrInvalid("purpose is required")
	}
	reservation.Purpose = strings.TrimSpace(purpose)
	reservation.Notes = strings.TrimSpace(notes)
	reservation.Revision++
	reservation.UpdatedAt = "sequence"
	event, err := s.newEvent(reservation.ID, "updated", reservation.Applicant, "reservation details updated")
	if err != nil {
		return model.Reservation{}, err
	}
	if err := s.store.UpdateReservationAndEvent(reservation, event); err != nil {
		return model.Reservation{}, err
	}
	return reservation, nil
}

func (s *Service) SearchReservations(query string, request model.PageRequest) (model.Page[model.Reservation], error) {
	return s.FindReservations(model.ReservationFilter{Query: query}, request)
}

func (s *Service) SelectReservation(id string) (model.ReservationDetail, error) {
	return s.GetReservationDetail(id)
}

func (s *Service) ChangeReservationScope(id string, scope string) (model.Reservation, error) {
	reservation, err := s.GetReservation(id)
	if err != nil {
		return model.Reservation{}, err
	}
	if reservation.Status != "pending" {
		return model.Reservation{}, errors.New("only pending reservations may change scope")
	}
	reservation.Scope = model.ValidateScope(scope)
	reservation.Revision++
	event, err := s.newEvent(reservation.ID, "scope_changed", reservation.Applicant, reservation.Scope)
	if err != nil {
		return model.Reservation{}, err
	}
	if err := s.store.UpdateReservationAndEvent(reservation, event); err != nil {
		return model.Reservation{}, err
	}
	return reservation, nil
}

func (s *Service) appendEvent(reservationID string, action string, actor string, detail string) error {
	event, err := s.newEvent(reservationID, action, actor, detail)
	if err != nil {
		return err
	}
	return s.store.PutAuditEvent(event)
}

func (s *Service) newEvent(reservationID string, action string, actor string, detail string) (model.AuditEvent, error) {
	id, sequence, err := s.store.AllocateID("EVT")
	if err != nil {
		return model.AuditEvent{}, err
	}
	return model.AuditEvent{ID: id, ReservationID: reservationID, Action: action, Actor: actor, Detail: detail, Sequence: sequence}, nil
}

func IsNotFound(err error) bool {
	return errors.Is(err, store.ErrNotFound)
}

func (s *Service) CheckReservationAccess(actor access.Actor, id string, action string) (access.Decision, error) {
	reservation, err := s.GetReservation(id)
	if err != nil {
		return access.Decision{}, err
	}
	switch action {
	case "view":
		return access.CanView(actor, reservation), nil
	case "edit":
		return access.CanEdit(actor, reservation), nil
	case "review":
		return access.CanReview(actor, reservation), nil
	case "archive":
		return access.CanArchive(actor, reservation), nil
	default:
		return access.Decision{Action: action, Reason: "unknown action"}, nil
	}
}

func (s *Service) ListForMember(member string, request model.PageRequest) (model.Page[model.Reservation], error) {
	actor := access.NewActor(member, access.RoleMember, "team")
	items, err := s.store.ListAllReservations(model.ReservationFilter{Scope: "all"})
	if err != nil {
		return model.Page[model.Reservation]{}, err
	}
	items = access.FilterVisible(actor, items)
	request = model.ValidatePage(request)
	start := (request.Page - 1) * request.PageSize
	if start > len(items) {
		start = len(items)
	}
	end := start + request.PageSize
	if end > len(items) {
		end = len(items)
	}
	return model.Page[model.Reservation]{Items: items[start:end], Page: request.Page, PageSize: request.PageSize, Total: len(items), HasNext: end < len(items)}, nil
}

func (s *Service) TimelineSummary(id string, scope string) (timeline.Summary, error) {
	detail, err := s.GetReservationDetail(id)
	if err != nil {
		return timeline.Summary{}, err
	}
	entries := make([]timeline.Entry, 0, len(detail.Timeline))
	for _, entry := range detail.Timeline {
		entries = append(entries, timeline.Entry{ID: entry.ID, ReservationID: id, Kind: entry.Kind, Actor: entry.Actor, Summary: entry.Summary, Sequence: entry.Sequence, Visibility: entry.Visibility})
	}
	return timeline.Summarize(timeline.Filter(entries, scope)), nil
}
