package service

import (
	"strings"

	"venue-reservation/internal/model"
)

func (s *Service) AddCollaboration(reservationID string, author string, body string, visibility string) (model.CollaborationNote, error) {
	reservation, err := s.GetReservation(reservationID)
	if err != nil {
		return model.CollaborationNote{}, err
	}
	author = strings.TrimSpace(author)
	body = strings.TrimSpace(body)
	visibility = model.ValidateScope(visibility)
	noteID, sequence, err := s.store.AllocateID("NOTE")
	if err != nil {
		return model.CollaborationNote{}, err
	}
	note := model.CollaborationNote{ID: noteID, ReservationID: reservation.ID, Author: author, Body: body, Visibility: visibility, Sequence: sequence}
	if err := model.ValidateCollaboration(note); err != nil {
		return model.CollaborationNote{}, err
	}
	if err := s.store.PutCollaborationNote(note); err != nil {
		return model.CollaborationNote{}, err
	}
	if err := s.appendEvent(reservation.ID, "collaborated", author, body); err != nil {
		return model.CollaborationNote{}, err
	}
	return note, nil
}

func (s *Service) ListCollaboration(reservationID string) ([]model.CollaborationNote, error) {
	return s.store.ListNotes(reservationID)
}

func (s *Service) VisibleNotes(reservationID string, scope string) ([]model.CollaborationNote, error) {
	notes, err := s.ListCollaboration(reservationID)
	if err != nil {
		return nil, err
	}
	scope = model.ValidateScope(scope)
	if scope == "all" {
		return notes, nil
	}
	visible := make([]model.CollaborationNote, 0, len(notes))
	for _, note := range notes {
		if note.Visibility == scope || note.Visibility == "all" {
			visible = append(visible, note)
		}
	}
	return visible, nil
}

func (s *Service) AddTeamNote(reservationID string, author string, body string) (model.CollaborationNote, error) {
	return s.AddCollaboration(reservationID, author, body, "team")
}

func (s *Service) AddPublicNote(reservationID string, author string, body string) (model.CollaborationNote, error) {
	return s.AddCollaboration(reservationID, author, body, "public")
}
