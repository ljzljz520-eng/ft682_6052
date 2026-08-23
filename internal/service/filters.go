package service

import (
	"strings"

	"venue-reservation/internal/model"
)

func FilterForMember(items []model.Reservation, member string) []model.Reservation {
	member = strings.TrimSpace(member)
	result := make([]model.Reservation, 0, len(items))
	for _, item := range items {
		if item.Applicant == member || item.Scope == "public" {
			result = append(result, item)
		}
	}
	return result
}

func FilterForReviewer(items []model.Reservation, status string) []model.Reservation {
	result := make([]model.Reservation, 0, len(items))
	for _, item := range items {
		if status == "" || item.Status == status {
			result = append(result, item)
		}
	}
	return result
}

func QueryFromValues(venueID string, applicant string, status string, scope string, query string) model.ReservationFilter {
	return model.ReservationFilter{VenueID: strings.TrimSpace(venueID), Applicant: strings.TrimSpace(applicant), Status: strings.TrimSpace(status), Scope: model.ValidateScope(scope), Query: strings.TrimSpace(query)}
}

func IsVisibleToMember(item model.Reservation, member string) bool {
	return item.Scope == "public" || strings.TrimSpace(item.Applicant) == strings.TrimSpace(member)
}
