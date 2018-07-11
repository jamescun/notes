package models

// User represents a single user object in the database.
type User struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password []byte `json:"-"`

	Timestamp
}

// RowID returns the unique ID of the object in the database.
func (u User) RowID() int64 {
	return u.ID
}

// Users contains one or more user objects from the database.
type Users []*User
