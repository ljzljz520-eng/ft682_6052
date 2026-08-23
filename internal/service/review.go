package service

import (
	"errors"
	"strings"

	"venue-reservation/internal/model"
)

func (s *Service) ReviewReservation(id string, status string, actor string, detail string) (model.Reservation, error) {
	reservation, err := s.GetReservation(id)
	if err != nil {
		return model.Reservation{}, err
	}
	if reservation.Status != "pending" {
		return model.Reservation{}, errors.New("only pending reservations may be reviewed")
	}
	if err := model.ValidateReview(status, detail); err != nil {
		return model.Reservation{}, err
	}
	actor = strings.TrimSpace(actor)
	if actor == "" {
		return model.Reservation{}, model.ErrInvalid("review actor is required")
	}
	reservation.Status = status
	reservation.Revision++
	reservation.UpdatedAt = "sequence"
	event, err := s.newEvent(reservation.ID, "reviewed", actor, strings.TrimSpace(detail))
	if err != nil {
		return model.Reservation{}, err
	}
	if err := s.store.UpdateReservationAndEvent(reservation, event); err != nil {
		return model.Reservation{}, err
	}
	if status == "approved" {
		if err := s.store.SetSlotState(reservation.SlotID, false); err != nil {
			return model.Reservation{}, err
		}
	}
	return reservation, nil
}

func (s *Service) ApproveReservation(id string, actor string) (model.Reservation, error) {
	return s.ReviewReservation(id, "approved", actor, "approved by reviewer")
}

func (s *Service) RejectReservation(id string, actor string, reason string) (model.Reservation, error) {
	return s.ReviewReservation(id, "rejected", actor, reason)
}

func (s *Service) CanReview(reservation model.Reservation) bool {
	return reservation.Status == "pending" && reservation.Revision > 0
}
