package service

import (
	"errors"
	"strings"

	"venue-reservation/internal/model"
	"venue-reservation/internal/store"
)

type Service struct {
	store *store.Store
}

func New(storeRef *store.Store) *Service {
	return &Service{store: storeRef}
}

func (s *Service) RegisterVenue(venue model.Venue) error {
	venue.Name = strings.TrimSpace(venue.Name)
	venue.Address = strings.TrimSpace(venue.Address)
	if err := model.ValidateVenue(venue); err != nil {
		return err
	}
	if _, err := s.store.GetVenue(venue.ID); err == nil {
		return errors.New("venue already exists")
	} else if !errors.Is(err, store.ErrNotFound) {
		return err
	}
	return s.store.PutVenue(venue)
}

func (s *Service) UpdateVenue(venue model.Venue) error {
	if _, err := s.store.GetVenue(venue.ID); err != nil {
		return err
	}
	venue.Name = strings.TrimSpace(venue.Name)
	venue.Address = strings.TrimSpace(venue.Address)
	if err := model.ValidateVenue(venue); err != nil {
		return err
	}
	return s.store.PutVenue(venue)
}

func (s *Service) GetVenue(id string) (model.Venue, error) {
	return s.store.GetVenue(strings.TrimSpace(id))
}

func (s *Service) ListVenues(onlyEnabled bool) ([]model.Venue, error) {
	venues, err := s.store.ListVenues()
	if err != nil {
		return nil, err
	}
	if !onlyEnabled {
		return venues, nil
	}
	filtered := make([]model.Venue, 0, len(venues))
	for _, venue := range venues {
		if venue.Enabled {
			filtered = append(filtered, venue)
		}
	}
	return filtered, nil
}

func (s *Service) AddTimeSlot(slot model.TimeSlot) error {
	if _, err := s.GetVenue(slot.VenueID); err != nil {
		return err
	}
	if slot.ID == "" {
		id, _, err := s.store.AllocateID("SLOT")
		if err != nil {
			return err
		}
		slot.ID = id
	}
	if !slot.Available {
		slot.Available = true
	}
	if slot.State == "" {
		slot.State = "open"
	}
	return s.store.PutTimeSlot(slot)
}

func (s *Service) GetTimeSlot(id string) (model.TimeSlot, error) {
	return s.store.GetTimeSlot(strings.TrimSpace(id))
}

func (s *Service) ListTimeSlots(venueID string, availableOnly bool) ([]model.TimeSlot, error) {
	slots, err := s.store.ListTimeSlots(strings.TrimSpace(venueID))
	if err != nil {
		return nil, err
	}
	if !availableOnly {
		return slots, nil
	}
	result := make([]model.TimeSlot, 0, len(slots))
	for _, slot := range slots {
		if slot.Available && slot.State == "open" {
			result = append(result, slot)
		}
	}
	return result, nil
}

func (s *Service) SetSlotAvailability(id string, available bool) error {
	slot, err := s.GetTimeSlot(id)
	if err != nil {
		return err
	}
	slot.Available = available
	if available {
		slot.State = "open"
	} else {
		slot.State = "held"
	}
	return s.store.PutTimeSlot(slot)
}
