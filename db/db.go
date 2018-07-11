package db

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"

	"github.com/lib/pq"
)

// ErrNotFound is returned when a query is expected to yield a single
// row, however no row matched the query.
var ErrNotFound = errors.New("record not found")

// ErrDuplicate is returned when a unique constraint on a column is
// not applicable.
var ErrDuplicate = errors.New("duplicate record")

// DB contains methods that are common to all database controllers.
type DB struct {
	*sql.DB
}

// New returns a new PostgreSQL database connection.
func New(dsn string) (*DB, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	return &DB{
		DB: db,
	}, nil
}

// Insert the case with PostgreSQL where the ID has to explicitly be returned
// and scanned into the ID field, where sql.Result LastInsertId is not usable.
func (d *DB) Insert(ctx context.Context, query string, id interface{}, args ...interface{}) error {
	return d.DB.QueryRowContext(ctx, query, args...).Scan(id)
}

// Tx executes the given function and database calls inside a transaction,
// committing the results on success, and rolling back on error/panic.
func (d *DB) Tx(ctx context.Context, fn func(*sql.Tx) error) (err error) {
	tx, err := d.DB.BeginTx(ctx, nil)
	if err != nil {
		return
	}
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p) // re-throw panic after Rollback
		} else if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()

	err = fn(tx)
	return err
}

// IsDuplicate extracts a unique_constraint_violation from a database error.
func IsDuplicate(err error) (constraint string, ok bool) {
	if err == nil {
		return
	}

	if pgErr, ok1 := err.(*pq.Error); ok1 {
		if pgErr.Severity == "ERROR" && pgErr.Code == "23505" {
			constraint = pgErr.Constraint
			ok = true
			return
		}
	}

	return
}

// Limit is a number which defines the maximum number of rows to return
// from a query with more than one result. If less than 1, all results
// will be returned.
type Limit int64

// Value implements database/sql/driver.Valuer
func (l Limit) Value() (driver.Value, error) {
	const all = "all"

	if l < 1 {
		return all, nil
	}

	return int64(l), nil
}

// DatabaseError is returned when an issue is encountered when executing a
// query or marshalling/unmarshalling the result.
type DatabaseError struct {
	msg   string
	cause error
}

// NewDatabaseError initialises a new DatabaseError with a message and
// optional cause. If the cause given is not nil and already of type
// DatabaseError, that will be returned instead.
func NewDatabaseError(msg string, cause error) *DatabaseError {
	if dbErr, ok := cause.(*DatabaseError); ok {
		return dbErr
	}

	return &DatabaseError{
		msg:   msg,
		cause: cause,
	}
}

// Cause returns the causal error, or nil if not configured.
func (d *DatabaseError) Cause() error {
	return d.cause
}

func (d *DatabaseError) Error() string {
	if d.cause != nil {
		return d.msg + ": " + d.cause.Error()
	}

	return d.msg
}
