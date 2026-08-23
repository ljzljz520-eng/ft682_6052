package access

import (
	"strings"

	"venue-reservation/internal/model"
)

type Role string

const (
	RoleMember   Role = "member"
	RoleReviewer Role = "reviewer"
	RoleOperator Role = "operator"
	RoleAdmin    Role = "admin"
)

type Actor struct {
	ID    string
	Role  Role
	Scope string
}

type Decision struct {
	Allowed bool
	Reason  string
	Action  string
}

func NewActor(id string, role Role, scope string) Actor {
	return Actor{ID: strings.TrimSpace(id), Role: NormalizeRole(role), Scope: model.ValidateScope(scope)}
}

func NormalizeRole(role Role) Role {
	switch role {
	case RoleReviewer, RoleOperator, RoleAdmin:
		return role
	default:
		return RoleMember
	}
}

func IsPrivileged(actor Actor) bool {
	return actor.Role == RoleReviewer || actor.Role == RoleOperator || actor.Role == RoleAdmin
}

func CanView(actor Actor, reservation model.Reservation) Decision {
	if actor.ID == "" {
		return denied("view", "actor is required")
	}
	if reservation.Scope == "public" || IsPrivileged(actor) || reservation.Applicant == actor.ID {
		return allowed("view", "reservation is visible in actor scope")
	}
	return denied("view", "reservation is outside actor scope")
}

func CanEdit(actor Actor, reservation model.Reservation) Decision {
	if decision := CanView(actor, reservation); !decision.Allowed {
		return decisionFor("edit", decision.Reason)
	}
	if reservation.Status != "pending" {
		return denied("edit", "only pending reservations can be edited")
	}
	if IsPrivileged(actor) || reservation.Applicant == actor.ID {
		return allowed("edit", "actor owns or administers the pending reservation")
	}
	return denied("edit", "actor does not own the pending reservation")
}

func CanReview(actor Actor, reservation model.Reservation) Decision {
	if actor.ID == "" {
		return denied("review", "actor is required")
	}
	if !IsPrivileged(actor) {
		return denied("review", "member cannot review reservations")
	}
	if reservation.Status != "pending" {
		return denied("review", "only pending reservations can be reviewed")
	}
	return allowed("review", "reviewer role is present")
}

func CanArchive(actor Actor, reservation model.Reservation) Decision {
	if actor.ID == "" {
		return denied("archive", "actor is required")
	}
	if !IsPrivileged(actor) {
		return denied("archive", "member cannot archive reservations")
	}
	if reservation.Status == "archived" {
		return denied("archive", "reservation is already archived")
	}
	return allowed("archive", "operator role is present")
}

func VisibleTo(actor Actor, reservation model.Reservation) bool {
	return CanView(actor, reservation).Allowed
}

func FilterVisible(actor Actor, reservations []model.Reservation) []model.Reservation {
	visible := make([]model.Reservation, 0, len(reservations))
	for _, reservation := range reservations {
		if VisibleTo(actor, reservation) {
			visible = append(visible, reservation)
		}
	}
	return visible
}

func ReasonFor(action string, actor Actor, reservation model.Reservation) string {
	decision := decisionFor(action, "")
	switch action {
	case "view":
		decision = CanView(actor, reservation)
	case "edit":
		decision = CanEdit(actor, reservation)
	case "review":
		decision = CanReview(actor, reservation)
	case "archive":
		decision = CanArchive(actor, reservation)
	}
	return decision.Reason
}

func allowed(action string, reason string) Decision {
	return Decision{Allowed: true, Action: action, Reason: reason}
}

func denied(action string, reason string) Decision {
	return Decision{Allowed: false, Action: action, Reason: reason}
}

func decisionFor(action string, reason string) Decision {
	return Decision{Action: action, Reason: reason}
}

func RoleCanPublish(role Role) bool {
	return role == RoleReviewer || role == RoleOperator || role == RoleAdmin
}

func RoleCanManageCatalog(role Role) bool {
	return role == RoleOperator || role == RoleAdmin
}

func RoleName(role Role) string {
	return string(NormalizeRole(role))
}

func ScopeName(actor Actor) string {
	if actor.Scope == "" {
		return "all"
	}
	return model.ValidateScope(actor.Scope)
}
