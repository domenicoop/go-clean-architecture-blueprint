package domain

import "time"

// Entity represents a generic domain entity.
type Entity struct {
	ID        string
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}
