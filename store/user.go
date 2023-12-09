package store

import (
	"context"
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
		return nil, err
	}
	return u, nil
}

// implement all other methods
func (s *pgUserStore) GetOne(ctx context.Context, id int) (User, error) {
	return User{}, nil
}

func (s *pgUserStore) GetOneByStudentID(ctx context.Context, studentID string) (User, error) {
	return User{}, nil
}

func (s *pgUserStore) GetOneByEmail(ctx context.Context, email string) (User, error) {
	return User{}, nil
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
