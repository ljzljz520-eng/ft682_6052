package service

import (
	"fmt"
	"strings"

	"venue-reservation/internal/model"
)

func (s *Service) ValidateImport(rows []model.ImportRow) []model.ImportIssue {
	issues := make([]model.ImportIssue, 0)
	seen := make(map[string]bool)
	for _, row := range rows {
		if strings.TrimSpace(row.Reference) == "" {
			issues = append(issues, model.ImportIssue{Reference: row.Reference, Message: "reference is required"})
			continue
		}
		if seen[row.Reference] {
			issues = append(issues, model.ImportIssue{Reference: row.Reference, Message: "reference is duplicated"})
			continue
		}
		seen[row.Reference] = true
		if _, err := s.GetVenue(row.VenueID); err != nil {
			issues = append(issues, model.ImportIssue{Reference: row.Reference, Message: "venue is unknown"})
			continue
		}
		if _, err := s.GetTimeSlot(row.SlotID); err != nil {
			issues = append(issues, model.ImportIssue{Reference: row.Reference, Message: "slot is unknown"})
			continue
		}
		if strings.TrimSpace(row.Applicant) == "" || strings.TrimSpace(row.Purpose) == "" {
			issues = append(issues, model.ImportIssue{Reference: row.Reference, Message: "applicant and purpose are required"})
		}
	}
	return issues
}

func (s *Service) ImportReservations(rows []model.ImportRow) (model.ImportReport, error) {
	issues := s.ValidateImport(rows)
	rejected := make(map[string]bool)
	for _, issue := range issues {
		rejected[issue.Reference] = true
	}
	report := model.ImportReport{Rejected: len(issues), Issues: issues, IDs: make([]string, 0)}
	for _, row := range rows {
		if rejected[row.Reference] {
			continue
		}
		reservation, err := s.SubmitReservation(model.Reservation{VenueID: row.VenueID, SlotID: row.SlotID, Applicant: row.Applicant, Purpose: row.Purpose, Scope: "team"})
		if err != nil {
			report.Rejected++
			report.Issues = append(report.Issues, model.ImportIssue{Reference: row.Reference, Message: err.Error()})
			continue
		}
		report.Accepted++
		report.IDs = append(report.IDs, reservation.ID)
	}
	return report, nil
}

func (s *Service) ImportSummary(report model.ImportReport) string {
	return strings.Join([]string{"accepted", formatCount(report.Accepted), "rejected", formatCount(report.Rejected)}, " ")
}

func formatCount(value int) string {
	if value < 0 {
		return "0"
	}
	if value < 10 {
		return string(rune('0' + value))
	}
	return fmt.Sprintf("%d", value)
}
