package timeline

import (
	"sort"
	"strings"

	"venue-reservation/internal/model"
)

type Entry struct {
	ID            string `json:"id"`
	ReservationID string `json:"reservation_id"`
	Kind          string `json:"kind"`
	Actor         string `json:"actor"`
	Summary       string `json:"summary"`
	Sequence      uint64 `json:"sequence"`
	Visibility    string `json:"visibility"`
}

type Summary struct {
	Entries int    `json:"entries"`
	Latest  string `json:"latest"`
	Actors  int    `json:"actors"`
}

func Build(reservation model.Reservation, events []model.AuditEvent, notes []model.CollaborationNote, archive *model.ArchiveRecord) []Entry {
	entries := make([]Entry, 0, len(events)+len(notes)+1)
	for _, event := range events {
		entries = append(entries, Entry{ID: event.ID, ReservationID: event.ReservationID, Kind: event.Action, Actor: event.Actor, Summary: event.Detail, Sequence: event.Sequence, Visibility: "all"})
	}
	for _, note := range notes {
		entries = append(entries, Entry{ID: note.ID, ReservationID: note.ReservationID, Kind: "note", Actor: note.Author, Summary: note.Body, Sequence: note.Sequence, Visibility: model.ValidateScope(note.Visibility)})
	}
	if archive != nil {
		entries = append(entries, Entry{ID: archive.ID, ReservationID: archive.ReservationID, Kind: "archive", Actor: archive.ArchivedBy, Summary: archive.Reason, Sequence: sequenceFromArchive(archive), Visibility: "all"})
	}
	sort.SliceStable(entries, func(i int, j int) bool { return entries[i].Sequence < entries[j].Sequence })
	return entries
}

func sequenceFromArchive(archive *model.ArchiveRecord) uint64 {
	if archive == nil {
		return 0
	}
	if strings.HasPrefix(archive.ID, "ARC-") {
		return parseSequence(archive.ID[4:])
	}
	return 0
}

func parseSequence(value string) uint64 {
	var result uint64
	for _, char := range value {
		if char < '0' || char > '9' {
			return 0
		}
		result = result*10 + uint64(char-'0')
	}
	return result
}

func Filter(entries []Entry, scope string) []Entry {
	scope = model.ValidateScope(scope)
	if scope == "all" {
		return append([]Entry(nil), entries...)
	}
	filtered := make([]Entry, 0, len(entries))
	for _, entry := range entries {
		if entry.Visibility == "all" || entry.Visibility == scope {
			filtered = append(filtered, entry)
		}
	}
	return filtered
}

func Summarize(entries []Entry) Summary {
	actors := make(map[string]bool)
	latest := ""
	for _, entry := range entries {
		if entry.Actor != "" {
			actors[entry.Actor] = true
		}
		if entry.Summary != "" {
			latest = entry.Summary
		}
	}
	return Summary{Entries: len(entries), Latest: latest, Actors: len(actors)}
}

func HasAction(entries []Entry, action string) bool {
	for _, entry := range entries {
		if entry.Kind == action {
			return true
		}
	}
	return false
}

func CountByKind(entries []Entry) map[string]int {
	counts := make(map[string]int)
	for _, entry := range entries {
		counts[entry.Kind]++
	}
	return counts
}

func Labels(entries []Entry) []string {
	labels := make([]string, 0, len(entries))
	for _, entry := range entries {
		labels = append(labels, Label(entry.Kind))
	}
	return labels
}

func Label(kind string) string {
	switch kind {
	case "submitted":
		return "Submitted"
	case "reviewed":
		return "Reviewed"
	case "updated":
		return "Updated"
	case "collaborated":
		return "Collaborated"
	case "archive":
		return "Archived"
	default:
		return "Activity"
	}
}

func Latest(entries []Entry) (Entry, bool) {
	if len(entries) == 0 {
		return Entry{}, false
	}
	result := entries[0]
	for _, entry := range entries[1:] {
		if entry.Sequence > result.Sequence {
			result = entry
		}
	}
	return result, true
}
