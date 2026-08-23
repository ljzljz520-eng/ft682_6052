package model

import "fmt"

type Venue struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Address  string `json:"address"`
	Capacity int    `json:"capacity"`
	Enabled  bool   `json:"enabled"`
}

type TimeSlot struct {
	ID        string `json:"id"`
	VenueID   string `json:"venue_id"`
	StartsAt  string `json:"starts_at"`
	EndsAt    string `json:"ends_at"`
	State     string `json:"state"`
	Available bool   `json:"available"`
}

type Reservation struct {
	ID        string `json:"id"`
	VenueID   string `json:"venue_id"`
	SlotID    string `json:"slot_id"`
	Applicant string `json:"applicant"`
	Purpose   string `json:"purpose"`
	Scope     string `json:"scope"`
	Status    string `json:"status"`
	Notes     string `json:"notes"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	Revision  int    `json:"revision"`
}

type AuditEvent struct {
	ID            string `json:"id"`
	ReservationID string `json:"reservation_id"`
	Action        string `json:"action"`
	Actor         string `json:"actor"`
	Detail        string `json:"detail"`
	Sequence      uint64 `json:"sequence"`
}

type CollaborationNote struct {
	ID            string `json:"id"`
	ReservationID string `json:"reservation_id"`
	Author        string `json:"author"`
	Body          string `json:"body"`
	Visibility    string `json:"visibility"`
	Sequence      uint64 `json:"sequence"`
}

type ArchiveRecord struct {
	ID            string `json:"id"`
	ReservationID string `json:"reservation_id"`
	Reason        string `json:"reason"`
	ArchivedBy    string `json:"archived_by"`
	ArchivedAt    string `json:"archived_at"`
}

type ReservationDetail struct {
	Reservation Reservation         `json:"reservation"`
	Venue       Venue               `json:"venue"`
	Slot        TimeSlot            `json:"slot"`
	Events      []AuditEvent        `json:"events"`
	Notes       []CollaborationNote `json:"notes"`
	Archive     *ArchiveRecord      `json:"archive,omitempty"`
	Timeline    []TimelineEntry     `json:"timeline"`
}

type TimelineEntry struct {
	ID         string `json:"id"`
	Kind       string `json:"kind"`
	Actor      string `json:"actor"`
	Summary    string `json:"summary"`
	Sequence   uint64 `json:"sequence"`
	Visibility string `json:"visibility"`
}

type ReservationFilter struct {
	VenueID   string
	Applicant string
	Status    string
	Scope     string
	Query     string
}

type PageRequest struct {
	Page     int
	PageSize int
}

type Page[T any] struct {
	Items    []T  `json:"items"`
	Page     int  `json:"page"`
	PageSize int  `json:"page_size"`
	Total    int  `json:"total"`
	HasNext  bool `json:"has_next"`
}

type ImportRow struct {
	Reference string `json:"reference"`
	VenueID   string `json:"venue_id"`
	SlotID    string `json:"slot_id"`
	Applicant string `json:"applicant"`
	Purpose   string `json:"purpose"`
}

type ImportIssue struct {
	Reference string `json:"reference"`
	Message   string `json:"message"`
}

type ImportReport struct {
	Accepted int           `json:"accepted"`
	Rejected int           `json:"rejected"`
	Issues   []ImportIssue `json:"issues"`
	IDs      []string      `json:"ids"`
}

type VenueReport struct {
	VenueID       string `json:"venue_id"`
	VenueName     string `json:"venue_name"`
	Total         int    `json:"total"`
	Pending       int    `json:"pending"`
	Approved      int    `json:"approved"`
	Rejected      int    `json:"rejected"`
	Archived      int    `json:"archived"`
	Collaboration int    `json:"collaboration_notes"`
}

func ReservationStatuses() []string {
	return []string{"pending", "approved", "rejected", "archived"}
}

func NewReservationID(sequence uint64) string {
	return fmt.Sprintf("RSV-%04d", sequence)
}

func NewEventID(sequence uint64) string {
	return fmt.Sprintf("EVT-%04d", sequence)
}

func NewNoteID(sequence uint64) string {
	return fmt.Sprintf("NOTE-%04d", sequence)
}

func NewArchiveID(sequence uint64) string {
	return fmt.Sprintf("ARC-%04d", sequence)
}

func NewTimeSlotID(venueID string, sequence uint64) string {
	return fmt.Sprintf("%s-SLOT-%02d", venueID, sequence)
}
