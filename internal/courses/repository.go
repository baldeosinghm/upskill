package courses

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Create a course

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
// COURSE FUNCTIONALITY
//

// Create a course
func (r *Repository) Create(ctx context.Context, name, ownerID string) (*Course, error) {
	var course Course
	err := r.db.QueryRow(
		ctx,
		`
		INSERT INTO courses (name, owner_id)
		VALUES($1, $2)
		RETURNING id, name, owner_id, created_at, updated_at
		`,
		name, ownerID,
	).Scan(&course.ID, &course.Name, &course.OwnerID, &course.CreatedAt, &course.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("create course: %w", err)
	}
	return &course, nil
}

// Get a course
func (r *Repository) GetByID(ctx context.Context, id string) (*Course, error) {
	var course Course
	err := r.db.QueryRow(
		ctx,
		`
		SELECT id, name, owner_id, created_at, updated_at FROM courses WHERE id=$1
		`,
		id,
	).Scan(&course.ID, &course.Name, &course.OwnerID, &course.CreatedAt, &course.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("get course: %w", err)
	}
	return &course, nil
}

// List all courses
func (r *Repository) List(ctx context.Context) ([]Course, error) {
	// TODO: pagniate rows
	rows, err := r.db.Query(ctx, `SELECT id, name, owner_id, created_at, updated_at FROM courses`)
	if err != nil {
		return nil, fmt.Errorf("list courses: %w", err)
	}
	defer rows.Close()

	var courses []Course
	for rows.Next() {
		var c Course
		if err := rows.Scan(&c.ID, &c.Name, &c.OwnerID, &c.CreatedAt, &c.UpdatedAt); err != nil {
			return nil, fmt.Errorf("list courses scan: %w", err)
		}
		courses = append(courses, c)
	}
	// If error within rows returned, log it
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("list courses iteration: %w", err)
	}
	return courses, nil
}
