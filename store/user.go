package store

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/jmoiron/sqlx"
)

type User struct {
	ID        int       `db:"id" json:"id"`
	StudentID string    `db:"student_id" json:"studentId"`
	Name      string    `db:"name" json:"name"`
	Email     string    `db:"email" json:"email"`
	CreatedAt time.Time `db:"created_at" json:"createdAt"`
}

type UserStore interface {
	GetAll(ctx context.Context) ([]User, error)
	GetOne(ctx context.Context, id int) (*User, error)
	GetOneByStudentID(ctx context.Context, studentID string) (*User, error)
	GetOneByEmail(ctx context.Context, email string) (*User, error)
	Create(ctx context.Context, user User) (*User, error)
	Update(ctx context.Context, user User) (*User, error)
	Delete(ctx context.Context, id int) error
}

type pgUserStore struct {
	db *sqlx.DB
}

func NewPGUserStore(db *sqlx.DB) UserStore {
	return &pgUserStore{db}
}

func (s *pgUserStore) GetAll(ctx context.Context) ([]User, error) {
	var u []User
	err := s.db.SelectContext(ctx, &u, "SELECT * FROM users")
	if err != nil {
		return nil, fmt.Errorf("Unable to retrieve users from database, error %w", err)
	}
	return u, nil
}

// implement all other methods
func (s *pgUserStore) GetOne(ctx context.Context, id int) (*User, error) {
	var u User
	err := s.db.GetContext(ctx, &u, "SELECT * FROM users WHERE id = $1", id)
	if err != nil {
		return nil, fmt.Errorf("Unable to find user with id, error %w", err)
	}
	return &u, nil
}

func (s *pgUserStore) GetOneByStudentID(ctx context.Context, studentID string) (*User, error) {
	var u User

	// Check if the length of the studentID is 8 characters
	if len(studentID) != 8 {
		return nil, fmt.Errorf("studentID must be 8 digits long, got %d", len(studentID))
	}

	// Convert string parameter to studentID integer type
	value, err := strconv.Atoi(studentID)
	if err != nil {
		return nil, fmt.Errorf("studentID must be a valid integer")
	}

	err = s.db.GetContext(ctx, &u, "SELECT * FROM users WHERE student_id = $1", value)
	if err != nil {
		return nil, fmt.Errorf("Unable to find user with studentID %d, with error %w", value, err)
	}
	return &u, nil
}

func (s *pgUserStore) GetOneByEmail(ctx context.Context, email string) (*User, error) {
	var u User

	err := s.db.GetContext(ctx, &u, "SELECT * FROM users WHERE TRIM(LOWER(email)) = TRIM(LOWER($1))", email)
	if err != nil {
		return nil, fmt.Errorf("Unable to find user with email %s, with error %w", email, err)
	}
	return &u, nil
}

func (s *pgUserStore) Create(ctx context.Context, user User) (*User, error) {
	// String StudentID conversion to integer
	value, err := strconv.Atoi(user.StudentID)
	if err != nil {
		return nil, fmt.Errorf("studentID must be a valid integer")
	}

	query := "INSERT INTO users (email, name, student_id) VALUES ($1, $2, $3) RETURNING id, created_at"
	// reassign the error value if needed
	err = s.db.QueryRowContext(ctx, query, user.Email, user.Name, value).Scan(&user.ID, &user.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("Unable to create and store a user, error %w", err)

	}
	return &user, nil
}

func (s *pgUserStore) Update(ctx context.Context, user User) (*User, error) {
	value, err := strconv.Atoi(user.StudentID)
	if err != nil {
		return nil, fmt.Errorf("studentID must be a valid integer")
	}

	query := "UPDATE users SET email = $1, name = $2, student_id = $3 WHERE id = $4"

	_, err := s.db.ExecContext(ctx, query, user.Email, user.Name, value, user.ID)

	if err != nil {
		return nil, fmt.Errorf("Unable to update user, error %w", err)

	}

	return &user, nil
}

func (s *pgUserStore) Delete(ctx context.Context, id int) error {
	query := "DELETE FROM users WHERE id = $1"

	_, err := s.db.ExecContext(ctx, query, id)

	if err != nil {
		return fmt.Errorf("Unable to delete user, error %w", err)

	}

	return nil
}
