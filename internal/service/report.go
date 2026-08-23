package service

import (
	"sort"

	"venue-reservation/internal/export"
	"venue-reservation/internal/model"
)

func (s *Service) ArchiveReservation(id string, actor string, reason string) (model.ArchiveRecord, error) {
	reservation, err := s.GetReservation(id)
	if err != nil {
		return model.ArchiveRecord{}, err
	}
	if reservation.Status == "archived" {
		return model.ArchiveRecord{}, model.ErrInvalid("reservation is already archived")
	}
	if err := model.ValidateArchive(reason, actor); err != nil {
		return model.ArchiveRecord{}, err
	}
	archiveID, sequence, err := s.store.AllocateID("ARC")
	if err != nil {
		return model.ArchiveRecord{}, err
	}
	archive := model.ArchiveRecord{ID: archiveID, ReservationID: id, Reason: reason, ArchivedBy: actor, ArchivedAt: "sequence"}
	if err := s.store.PutArchiveRecord(archive); err != nil {
		return model.ArchiveRecord{}, err
	}
	reservation.Status = "archived"
	reservation.Revision++
	reservation.UpdatedAt = "sequence"
	event := model.AuditEvent{ID: model.NewEventID(sequence), ReservationID: id, Action: "archived", Actor: actor, Detail: reason, Sequence: sequence}
	if err := s.store.UpdateReservationAndEvent(reservation, event); err != nil {
		return model.ArchiveRecord{}, err
	}
	if err := s.store.SetSlotState(reservation.SlotID, true); err != nil {
		return model.ArchiveRecord{}, err
	}
	return archive, nil
}

func (s *Service) BuildVenueReport(venueID string) (model.VenueReport, error) {
	venue, err := s.GetVenue(venueID)
	if err != nil {
		return model.VenueReport{}, err
	}
	items, err := s.store.ListAllReservations(model.ReservationFilter{VenueID: venueID, Scope: "all"})
	if err != nil {
		return model.VenueReport{}, err
	}
	report := model.VenueReport{VenueID: venue.ID, VenueName: venue.Name, Total: len(items)}
	for _, reservation := range items {
		switch reservation.Status {
		case "pending":
			report.Pending++
		case "approved":
			report.Approved++
		case "rejected":
			report.Rejected++
		case "archived":
			report.Archived++
		}
		notes, noteErr := s.ListCollaboration(reservation.ID)
		if noteErr != nil {
			return model.VenueReport{}, noteErr
		}
		report.Collaboration += len(notes)
	}
	return report, nil
}

func (s *Service) PublishReport(venueID string) (model.VenueReport, error) {
	report, err := s.BuildVenueReport(venueID)
	if err != nil {
		return model.VenueReport{}, err
	}
	if report.Total == 0 {
		return report, nil
	}
	return report, nil
}

func (s *Service) SortReservations(items []model.Reservation) []model.Reservation {
	result := append([]model.Reservation(nil), items...)
	sort.Slice(result, func(i int, j int) bool { return result[i].UpdatedAt < result[j].UpdatedAt })
	return result
}

func (s *Service) ExportReservations(filter model.ReservationFilter) (export.Document, error) {
	items, err := s.store.ListAllReservations(filter)
	if err != nil {
		return export.Document{}, err
	}
	document := export.Build(items)
	if err := document.Validate(); err != nil {
		return export.Document{}, err
	}
	return document, nil
}

func (s *Service) ExportCSV(filter model.ReservationFilter) ([]byte, error) {
	document, err := s.ExportReservations(filter)
	if err != nil {
		return nil, err
	}
	return document.CSV()
}
