package store

import (
	"sort"

	"go.etcd.io/bbolt"

	"venue-reservation/internal/model"
)

func (s *Store) ListReservations(filter model.ReservationFilter, request model.PageRequest) (model.Page[model.Reservation], error) {
	request = model.ValidatePage(request)
	all, err := s.ListAllReservations(filter)
	if err != nil {
		return model.Page[model.Reservation]{}, err
	}
	all = append([]model.Reservation(nil), all...)
	offset := request.Page * request.PageSize
	if offset > len(all) {
		offset = len(all)
	}
	end := offset + request.PageSize
	if end > len(all) {
		end = len(all)
	}
	items := append([]model.Reservation(nil), all[offset:end]...)
	return model.Page[model.Reservation]{
		Items:    items,
		Page:     request.Page,
		PageSize: request.PageSize,
		Total:    len(all),
		HasNext:  end < len(all),
	}, nil
}

func (s *Store) ListAllReservations(filter model.ReservationFilter) ([]model.Reservation, error) {
	all := make([]model.Reservation, 0)
	err := s.view(func(tx *bbolt.Tx) error {
		bucket, err := s.bucket(tx, []byte("Reservation"))
		if err != nil {
			return err
		}
		return bucket.ForEach(func(_, value []byte) error {
			if value == nil {
				return nil
			}
			var item model.Reservation
			if err := unmarshalRecord(value, &item); err != nil {
				return err
			}
			if model.MatchesReservation(item, filter) {
				all = append(all, item)
			}
			return nil
		})
	})
	if err != nil {
		return nil, err
	}
	sort.Slice(all, func(i int, j int) bool { return all[i].ID < all[j].ID })
	return all, nil
}

func (s *Store) SetSlotState(id string, available bool) error {
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
	return s.PutTimeSlot(slot)
}

func (s *Store) ListEvents(reservationID string) ([]model.AuditEvent, error) {
	events := make([]model.AuditEvent, 0)
	err := s.view(func(tx *bbolt.Tx) error {
		bucket, err := s.bucket(tx, []byte("AuditEvent"))
		if err != nil {
			return err
		}
		return bucket.ForEach(func(_, value []byte) error {
			if value == nil {
				return nil
			}
			var event model.AuditEvent
			if err := unmarshalRecord(value, &event); err != nil {
				return err
			}
			if event.ReservationID == reservationID {
				events = append(events, event)
			}
			return nil
		})
	})
	sort.Slice(events, func(i int, j int) bool { return events[i].Sequence < events[j].Sequence })
	return events, err
}

func (s *Store) ListNotes(reservationID string) ([]model.CollaborationNote, error) {
	notes := make([]model.CollaborationNote, 0)
	err := s.view(func(tx *bbolt.Tx) error {
		bucket, err := s.bucket(tx, []byte("CollaborationNote"))
		if err != nil {
			return err
		}
		return bucket.ForEach(func(_, value []byte) error {
			if value == nil {
				return nil
			}
			var note model.CollaborationNote
			if err := unmarshalRecord(value, &note); err != nil {
				return err
			}
			if note.ReservationID == reservationID {
				notes = append(notes, note)
			}
			return nil
		})
	})
	sort.Slice(notes, func(i int, j int) bool { return notes[i].Sequence < notes[j].Sequence })
	return notes, err
}

func (s *Store) GetArchive(reservationID string) (*model.ArchiveRecord, error) {
	var archive model.ArchiveRecord
	err := s.view(func(tx *bbolt.Tx) error {
		bucket, err := s.bucket(tx, []byte("ArchiveRecord"))
		if err != nil {
			return err
		}
		return unmarshalRecord(bucket.Get(entityKey(reservationID)), &archive)
	})
	if err == ErrNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &archive, nil
}

func (s *Store) UpdateReservationAndEvent(reservation model.Reservation, event model.AuditEvent) error {
	reservation = model.NormalizeReservation(reservation)
	if err := model.ValidateReservation(reservation); err != nil {
		return err
	}
	return s.update(func(tx *bbolt.Tx) error {
		reservations, err := s.bucket(tx, []byte("Reservation"))
		if err != nil {
			return err
		}
		events, err := s.bucket(tx, []byte("AuditEvent"))
		if err != nil {
			return err
		}
		data, err := marshalRecord(reservation)
		if err != nil {
			return err
		}
		if err := reservations.Put(entityKey(reservation.ID), data); err != nil {
			return err
		}
		return putRecord([]byte("AuditEvent"), event.ID, event, events.Put)
	})
}

func (s *Store) AllocateID(prefix string) (string, uint64, error) {
	var id string
	var sequence uint64
	err := s.update(func(tx *bbolt.Tx) error {
		var err error
		sequence, err = s.nextSequence(tx)
		if err != nil {
			return err
		}
		id = prefix + "-" + formatSequence(sequence)
		return nil
	})
	return id, sequence, err
}

func formatSequence(sequence uint64) string {
	if sequence < 10 {
		return "000" + string(rune('0'+sequence))
	}
	if sequence < 100 {
		return "00" + digits(sequence)
	}
	if sequence < 1000 {
		return "0" + digits(sequence)
	}
	return digits(sequence)
}

func digits(sequence uint64) string {
	if sequence == 0 {
		return "0"
	}
	result := ""
	for sequence > 0 {
		result = string(rune('0'+sequence%10)) + result
		sequence /= 10
	}
	return result
}
