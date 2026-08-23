package export

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"venue-reservation/internal/model"
)

type Row struct {
	ID        string
	VenueID   string
	SlotID    string
	Applicant string
	Purpose   string
	Scope     string
	Status    string
	Revision  int
}

type Document struct {
	Title   string   `json:"title"`
	Columns []string `json:"columns"`
	Rows    []Row    `json:"rows"`
}

func Build(items []model.Reservation) Document {
	rows := make([]Row, 0, len(items))
	for _, item := range items {
		rows = append(rows, Row{ID: item.ID, VenueID: item.VenueID, SlotID: item.SlotID, Applicant: item.Applicant, Purpose: item.Purpose, Scope: item.Scope, Status: item.Status, Revision: item.Revision})
	}
	sort.Slice(rows, func(i int, j int) bool { return rows[i].ID < rows[j].ID })
	return Document{Title: "venue reservation export", Columns: Columns(), Rows: rows}
}

func Columns() []string {
	return []string{"id", "venue_id", "slot_id", "applicant", "purpose", "scope", "status", "revision"}
}

func (document Document) Validate() error {
	if strings.TrimSpace(document.Title) == "" {
		return fmt.Errorf("export title is required")
	}
	if len(document.Columns) != len(Columns()) {
		return fmt.Errorf("export columns are incomplete")
	}
	if len(document.Rows) == 0 {
		return nil
	}
	for _, row := range document.Rows {
		if row.ID == "" || row.VenueID == "" || row.SlotID == "" {
			return fmt.Errorf("export row identity is incomplete")
		}
		if row.Revision < 1 {
			return fmt.Errorf("export row revision is invalid")
		}
	}
	return nil
}

func (document Document) JSON() ([]byte, error) {
	if err := document.Validate(); err != nil {
		return nil, err
	}
	return json.MarshalIndent(document, "", "  ")
}

func (document Document) CSV() ([]byte, error) {
	if err := document.Validate(); err != nil {
		return nil, err
	}
	var buffer bytes.Buffer
	writer := csv.NewWriter(&buffer)
	if err := writer.Write(document.Columns); err != nil {
		return nil, err
	}
	for _, row := range document.Rows {
		if err := writer.Write(rowValues(row)); err != nil {
			return nil, err
		}
	}
	writer.Flush()
	if err := writer.Error(); err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

func rowValues(row Row) []string {
	return []string{row.ID, row.VenueID, row.SlotID, row.Applicant, row.Purpose, row.Scope, row.Status, fmt.Sprintf("%d", row.Revision)}
}

func Filter(items []model.Reservation, filter model.ReservationFilter) []model.Reservation {
	result := make([]model.Reservation, 0, len(items))
	for _, item := range items {
		if model.MatchesReservation(item, filter) {
			result = append(result, item)
		}
	}
	return result
}

func StatusLabel(status string) string {
	switch status {
	case "pending":
		return "Pending review"
	case "approved":
		return "Approved"
	case "rejected":
		return "Rejected"
	case "archived":
		return "Archived"
	default:
		return "Unknown"
	}
}

func StatusLabels(items []model.Reservation) map[string]string {
	labels := make(map[string]string)
	for _, item := range items {
		labels[item.ID] = StatusLabel(item.Status)
	}
	return labels
}

func CountByStatus(items []model.Reservation) map[string]int {
	counts := make(map[string]int)
	for _, item := range items {
		counts[item.Status]++
	}
	return counts
}

func EmptyDocument() Document {
	return Document{Title: "venue reservation export", Columns: Columns(), Rows: []Row{}}
}
