package models

import (
	"time"
)

// Model is a common interface implemented by all objects capable of being
// serialized into the database.
type Model interface {
	// RowID returns the unique ID of the object in the database.
	RowID() int64

	// Created sets the CreatedAt and UpdatedAt timestamp of the object,
	// overwriting them if they have already been set, and resetting the
	// DeletedAt timestamp if set.
	Created()

	// Updated sets the UpdatedAt timestamp of the object, resetting the
	// DeletedAt timestamp if set.
	Updated()

	// Deleted sets the DeletedAt timestamp.
	Deleted()
}

// Timestamp implements the timestamps required of a database model.
type Timestamp struct {
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}

func (t *Timestamp) Created() {
	now := time.Now().UTC()
	t.CreatedAt = now
	t.UpdatedAt = now
	t.DeletedAt = nil
}

func (t *Timestamp) Updated() {
	now := time.Now().UTC()
	t.UpdatedAt = now
	t.DeletedAt = nil
}

func (t *Timestamp) Deleted() {
	now := time.Now().UTC()
	t.DeletedAt = &now
}
