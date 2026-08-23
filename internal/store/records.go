package store

import (
	"fmt"

	"go.etcd.io/bbolt"

	"venue-reservation/internal/model"
)

func (s *Store) PutVenue(venue model.Venue) error {
	if err := ensureVenue(venue); err != nil {
		return err
	}
	return s.update(func(tx *bbolt.Tx) error {
		bucket, err := s.bucket(tx, []byte("Venue"))
		if err != nil {
			return err
		}
		return putRecord([]byte("Venue"), venue.ID, venue, bucket.Put)
	})
}

func (s *Store) GetVenue(id string) (model.Venue, error) {
	var venue model.Venue
	err := s.view(func(tx *bbolt.Tx) error {
		bucket, err := s.bucket(tx, []byte("Venue"))
		if err != nil {
			return err
		}
		return unmarshalRecord(bucket.Get(entityKey(id)), &venue)
	})
	return venue, err
}

func (s *Store) ListVenues() ([]model.Venue, error) {
	venues := make([]model.Venue, 0)
	err := s.view(func(tx *bbolt.Tx) error {
		bucket, err := s.bucket(tx, []byte("Venue"))
		if err != nil {
			return err
		}
		return bucket.ForEach(func(_, value []byte) error {
			if value == nil {
				return nil
			}
			var venue model.Venue
			if err := unmarshalRecord(value, &venue); err != nil {
				return err
			}
			venues = append(venues, venue)
			return nil
		})
	})
	return venues, err
}

func (s *Store) PutTimeSlot(slot model.TimeSlot) error {
	if err := model.ValidateTimeSlot(slot); err != nil {
		return err
	}
	return s.update(func(tx *bbolt.Tx) error {
		bucket, err := s.bucket(tx, []byte("TimeSlot"))
		if err != nil {
			return err
		}
		return putRecord([]byte("TimeSlot"), slot.ID, slot, bucket.Put)
	})
}

func (s *Store) GetTimeSlot(id string) (model.TimeSlot, error) {
	var slot model.TimeSlot
	err := s.view(func(tx *bbolt.Tx) error {
		bucket, err := s.bucket(tx, []byte("TimeSlot"))
		if err != nil {
			return err
		}
		return unmarshalRecord(bucket.Get(entityKey(id)), &slot)
	})
	return slot, err
}

func (s *Store) ListTimeSlots(venueID string) ([]model.TimeSlot, error) {
	slots := make([]model.TimeSlot, 0)
	err := s.view(func(tx *bbolt.Tx) error {
		bucket, err := s.bucket(tx, []byte("TimeSlot"))
		if err != nil {
			return err
		}
		return bucket.ForEach(func(_, value []byte) error {
			if value == nil {
				return nil
			}
			var slot model.TimeSlot
			if err := unmarshalRecord(value, &slot); err != nil {
				return err
			}
			if venueID == "" || slot.VenueID == venueID {
				slots = append(slots, slot)
			}
			return nil
		})
	})
	return slots, err
}

func (s *Store) PutReservation(reservation model.Reservation) error {
	reservation = model.NormalizeReservation(reservation)
	if err := model.ValidateReservation(reservation); err != nil {
		return err
	}
	return s.update(func(tx *bbolt.Tx) error {
		bucket, err := s.bucket(tx, []byte("Reservation"))
		if err != nil {
			return err
		}
		return putRecord([]byte("Reservation"), reservation.ID, reservation, bucket.Put)
	})
}

func (s *Store) GetReservation(id string) (model.Reservation, error) {
	var reservation model.Reservation
	err := s.view(func(tx *bbolt.Tx) error {
		bucket, err := s.bucket(tx, []byte("Reservation"))
		if err != nil {
			return err
		}
		return unmarshalRecord(bucket.Get(entityKey(id)), &reservation)
	})
	return reservation, err
}

func (s *Store) deleteReservation(id string) error {
	return s.update(func(tx *bbolt.Tx) error {
		bucket, err := s.bucket(tx, []byte("Reservation"))
		if err != nil {
			return err
		}
		return bucket.Delete(entityKey(id))
	})
}

func (s *Store) PutAuditEvent(event model.AuditEvent) error {
	return s.update(func(tx *bbolt.Tx) error {
		bucket, err := s.bucket(tx, []byte("AuditEvent"))
		if err != nil {
			return err
		}
		return putRecord([]byte("AuditEvent"), event.ID, event, bucket.Put)
	})
}

func (s *Store) PutCollaborationNote(note model.CollaborationNote) error {
	return s.update(func(tx *bbolt.Tx) error {
		bucket, err := s.bucket(tx, []byte("CollaborationNote"))
		if err != nil {
			return err
		}
		return putRecord([]byte("CollaborationNote"), note.ID, note, bucket.Put)
	})
}

func (s *Store) PutArchiveRecord(archive model.ArchiveRecord) error {
	return s.update(func(tx *bbolt.Tx) error {
		bucket, err := s.bucket(tx, []byte("ArchiveRecord"))
		if err != nil {
			return err
		}
		return putRecord([]byte("ArchiveRecord"), archive.ReservationID, archive, bucket.Put)
	})
}

func combineWriteError(action string, err error) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("%s: %w", action, err)
}
