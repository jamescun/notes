package models

// Note represents a single Note object in the database.
type Note struct {
	ID     int64 `json:"id"`
	UserID int64 `json:"user_id"`

	Title *string `json:"title"`
	Body  string  `json:"body"`

	Timestamp
}

// RowID returns the unique ID of the object in the database.
func (n Note) RowID() int64 {
	return n.ID
}

// Notes contains one or more note objects from the database.
type Notes []*Note
