package users

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

// The repository file reads data then converts it into SQL statements
// that can be understood by our database.

// This struct wraps the db connection pool. It's the repository's
// only state - a reference to the database. All methods on this struct
// will use r.db to run queries. The repository should never create
// its own database connection.
type Repository struct {
	db *pgxpool.Pool
}

// This constructor allows us to create an instance of the
// Repository struct.  We'll use dependency injection (see main.go)
// to avoid testing db connection here, but rather from the main file.
func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

//
// TEACHER FUNCTIONALITY
//

// Create user: either TEACHER or STUDENT
func (r Repository) Create(ctx context.Context, username, email, passwordHash, role string) (*User, error) {
	// create var user of type User so that row data can be store in it
	var user User
	args := []any{username, email, passwordHash, role}
	err := r.db.QueryRow(
		ctx,
		`
		INSERT INTO users (username, email, password, role)
		VALUES($1, $2, $3, $4)
		RETURNING id, username, email, role, created_at
		`,
		args...,
	).Scan(&user.ID, &user.Username, &user.Email, &user.Role, &user.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}
	return &user, nil
}

// Locate user by email
func (r *Repository) FindByEmail(ctx context.Context, email string) (*User, error) {
	return nil, nil
}

// Delete user

// Login

// Create a course (make them owner)

// Create an assignment

// Invite student to course

// Share assignment w/ student

// Update assignment

// Delete assignment

// Grade assignment

//
// STUDENT FUNCTIONALITY
//

// Join a course (student)

// See assignments

// Turn in assignment (either no action or action, i.e. upload)
