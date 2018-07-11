package db

import (
	"context"
	"database/sql"

	"github.com/jamescun/notes/models"
)

// Users contains methods for interacting with user objects in the database.
type Users struct {
	*DB
}

// Create inserts a new user object in the database, setting the ID and
// timestamps in the given user struct.
func (u *Users) Create(ctx context.Context, x *models.User) error {
	const query = `
		INSERT INTO users
			(username, email, password, created_at, updated_at)
		VALUES
			($1, $2, $3, $4, $5)
		RETURNING id;
	`

	err := u.Insert(
		ctx, query, &x.ID, &x.Username, &x.Email,
		&x.Password, &x.CreatedAt, &x.UpdatedAt,
	)
	if _, ok := IsDuplicate(err); ok {
		return ErrDuplicate
	} else if err != nil {
		return NewDatabaseError("could not create user", err)
	}

	return nil
}

func (u *Users) get(ctx context.Context, id *int64, username, email *string) (*models.User, error) {
	const query = `
		SELECT
			id, username, email, password, created_at, updated_at, deleted_at
		FROM users
		WHERE
			id = $1 OR
			username = $2 OR
			email = $3
		LIMIT 1;
	`

	var user models.User

	err := u.QueryRowContext(ctx, query, id, username, email).Scan(
		&user.ID, &user.Username, &user.Email, &user.Password,
		&user.CreatedAt, &user.UpdatedAt, &user.DeletedAt,
	)
	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	} else if err != nil {
		return nil, NewDatabaseError("could not get user", err)
	}

	return &user, nil
}

// GetByID retrieves a single user from the database by their ID.
func (u *Users) GetByID(ctx context.Context, id int64) (*models.User, error) {
	return u.get(ctx, &id, nil, nil)
}

// GetByUsername retrieves a single user from the database by their Username.
func (u *Users) GetByUsername(ctx context.Context, username string) (*models.User, error) {
	return u.get(ctx, nil, &username, nil)
}

// GetByEmail retrieves a single user from the database by their Email Address.
func (u *Users) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	return u.get(ctx, nil, nil, &email)
}
