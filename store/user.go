package store

import (
	"context"
	"time"
	"strconv"
	"fmt"
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
	GetOne(ctx context.Context, id int) (User, error)
	GetOneByStudentID(ctx context.Context, studentID string) (User, error)
	GetOneByEmail(ctx context.Context, email string) (User, error)
	Create(ctx context.Context, user User) (User, error)
	Update(ctx context.Context, user User) (User, error)
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
		return nil,fmt.Errorf("Unable to retrieve users from database")
	}
	return u, nil
}

// implement all other methods
func (s *pgUserStore) GetOne(ctx context.Context, id int) (User, error) {
	var u User 
	err := s.db.SelectContext(ctx, &u, "SELECT * FROM users WHERE id = $1", id)
	if err != nil {
		
		return User{}, fmt.Errorf("Unable to find user with id")
	}
	return u, nil 
}

func (s *pgUserStore) GetOneByStudentID(ctx context.Context, studentID string) (User, error) {
    var u User

    // Check if the length of the studentID is 8 characters
    if len(studentID) != 8 {
        return User{}, fmt.Errorf("studentID must be 8 digits long, got %d", len(studentID))
    }

    // Convert string parameter to studentID integer type 
    value, err := strconv.Atoi(studentID)
    if err != nil {
        return User{}, fmt.Errorf("studentID must be a valid integer")
    }

    err = s.db.GetContext(ctx, &u, "SELECT * FROM users WHERE student_id = $1", value)
    if err != nil {
        return User{}, fmt.Errorf("Unable to find user with studentID %d", value)
    }
    return u, nil 
}


func (s *pgUserStore) GetOneByEmail(ctx context.Context, email string) (User, error) {
	var u User 
	
	err := s.db.GetContext(ctx, &u, "SELECT * FROM users WHERE TRIM(LOWER(email)) = TRIM(LOWER($1))", email)
	if err != nil {
		return User{}, fmt.Errorf("Unable to find user with email %s", email)
	}
	return u, nil
}

func (s *pgUserStore) Create(ctx context.Context, user User) (User, error) {
	return User{}, nil
}

func (s *pgUserStore) Update(ctx context.Context, user User) (User, error) {
	return User{}, nil
}

func (s *pgUserStore) Delete(ctx context.Context, id int) error {
	return nil
}
