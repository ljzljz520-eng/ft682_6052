package seed

import (
	"venue-reservation/internal/model"
	"venue-reservation/internal/service"
)

func EnsureDemoData(app *service.Service) error {
	venues, err := app.ListVenues(false)
	if err != nil {
		return err
	}
	if len(venues) > 0 {
		return nil
	}
	for _, venue := range catalog() {
		if err := app.RegisterVenue(venue); err != nil {
			return err
		}
	}
	for _, slot := range slots() {
		if err := app.AddTimeSlot(slot); err != nil {
			return err
		}
	}
	return nil
}

func catalog() []model.Venue {
	return []model.Venue{
		{ID: "A-101", Name: "North Workshop", Address: "Building A", Capacity: 40, Enabled: true},
		{ID: "B-201", Name: "South Forum", Address: "Building B", Capacity: 120, Enabled: true},
	}
}

func slots() []model.TimeSlot {
	return []model.TimeSlot{
		{ID: model.NewTimeSlotID("A-101", 1), VenueID: "A-101", StartsAt: "09:00", EndsAt: "10:00", State: "open", Available: true},
		{ID: model.NewTimeSlotID("A-101", 2), VenueID: "A-101", StartsAt: "10:00", EndsAt: "11:00", State: "open", Available: true},
		{ID: model.NewTimeSlotID("B-201", 1), VenueID: "B-201", StartsAt: "14:00", EndsAt: "15:30", State: "open", Available: true},
	}
}

func CatalogVenueIDs() []string {
	items := catalog()
	ids := make([]string, 0, len(items))
	for _, item := range items {
		ids = append(ids, item.ID)
	}
	return ids
}

func CatalogSlotsFor(venueID string) []model.TimeSlot {
	result := make([]model.TimeSlot, 0)
	for _, slot := range slots() {
		if slot.VenueID == venueID {
			result = append(result, slot)
		}
	}
	return result
}
