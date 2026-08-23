package model

import (
	"net/url"
	"strings"
)

func QueryFilter(values url.Values) ReservationFilter {
	return ReservationFilter{
		VenueID:   strings.TrimSpace(values.Get("venue_id")),
		Applicant: strings.TrimSpace(values.Get("applicant")),
		Status:    strings.TrimSpace(values.Get("status")),
		Scope:     ValidateScope(values.Get("scope")),
		Query:     strings.TrimSpace(values.Get("q")),
	}
}

func ContainsText(value string, query string) bool {
	needle := strings.ToLower(strings.TrimSpace(query))
	if needle == "" {
		return true
	}
	return strings.Contains(strings.ToLower(value), needle)
}

func MatchesReservation(reservation Reservation, filter ReservationFilter) bool {
	if filter.VenueID != "" && reservation.VenueID != filter.VenueID {
		return false
	}
	if filter.Applicant != "" && reservation.Applicant != filter.Applicant {
		return false
	}
	if filter.Status != "" && reservation.Status != filter.Status {
		return false
	}
	if filter.Scope != "" && filter.Scope != "all" && reservation.Scope != filter.Scope {
		return false
	}
	return ContainsText(reservation.ID, filter.Query) || ContainsText(reservation.Purpose, filter.Query)
}

func NormalizeReservation(reservation Reservation) Reservation {
	if reservation.Status == "" {
		reservation.Status = "pending"
	}
	if reservation.Scope == "" {
		reservation.Scope = "team"
	}
	if reservation.Revision < 1 {
		reservation.Revision = 1
	}
	return reservation
}

func StatusIsTerminal(status string) bool {
	return status == "rejected" || status == "archived"
}
