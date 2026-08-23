package access

import (
	"testing"

	"venue-reservation/internal/model"
)

func TestReservationAccessDecisions(t *testing.T) {
	reservation := model.Reservation{ID: "RSV-0001", Applicant: "member-a", Scope: "team", Status: "pending", Revision: 1}
	owner := NewActor("member-a", RoleMember, "team")
	reviewer := NewActor("reviewer", RoleReviewer, "all")
	if !CanView(owner, reservation).Allowed || !CanEdit(owner, reservation).Allowed {
		t.Fatal("owner should view and edit pending reservation")
	}
	if !CanReview(reviewer, reservation).Allowed || !CanArchive(reviewer, reservation).Allowed {
		t.Fatal("reviewer should review and archive")
	}
	if CanReview(owner, reservation).Allowed {
		t.Fatal("member should not review")
	}
}
