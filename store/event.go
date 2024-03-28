package store

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
)

type Event struct {
	ID                 int       `db:"id"`
	Name               string    `db:"name"`
	Description        string    `db:"description"`
	StartDate          time.Time `db:"start_date"`
	EndDate            time.Time `db:"end_date"`
	VolunteersRequired int       `db:"volunteers_required"`
	Location           string    `db:"location"`
	CreatedAt          time.Time `db:"created_at"`
}

type EventStore interface {
	GetAll(ctx context.Context) ([]Event, error)
	GetOne(ctx context.Context, id int) (*Event, error)
	GetOneByStartDate(ctx context.Context, startDate time.Time) (*Event, error)
	GetOneByName(ctx context.Context, name string) (*Event, error) // ??
	Create(ctx context.Context, name string, description string, start time.Time, end time.Time, volunteers int, location string) (*Event, error)
	Update(ctx context.Context, event Event) (*Event, error)
	Delete(ctx context.Context, id int) error
	UnregisterEvent(ctx context.Context, UserID int, EventID int) error
}

type pgEventStore struct {
	db *sqlx.DB
}

func NewPGEventStore(db *sqlx.DB) EventStore {
	return &pgEventStore{db}
}

func (s *pgEventStore) GetAll(ctx context.Context) ([]Event, error) {
	var e []Event
	err := s.db.SelectContext(ctx, &e, "SELECT * FROM events")
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve events from database, error %w", err)
	}
	return e, nil
}

func (s *pgEventStore) GetOne(ctx context.Context, id int) (*Event, error) {
	var e Event
	err := s.db.GetContext(ctx, &e, "SELECT * FROM events WHERE id = $1", id)
	if err != nil {
		return nil, fmt.Errorf("unable to find event with id, error %w", err)
	}
	return &e, nil
}

func (s *pgEventStore) GetOneByName(ctx context.Context, name string) (*Event, error) {
	var e Event
	err := s.db.GetContext(ctx, &e, "SELECT * FROM events WHERE LOWER(name) = LOWER($1)", name)
	if err != nil {
		return nil, fmt.Errorf("unable to find event with name %s, with error %w", name, err)
	}
	return &e, nil
}

// Don't know if I handled date right here
func (s *pgEventStore) GetOneByStartDate(ctx context.Context, date time.Time) (*Event, error) {
	var e Event
	formattedDate := date.Format("2006-01-02 15:04")
	err := s.db.GetContext(ctx, &e, "SELECT * FROM events WHERE TO_CHAR(start_date, 'YYYY-MM-DD HH24:MI') = $1", formattedDate)
	if err != nil {
		return nil, fmt.Errorf("unable to find event with start date %s, with error %w", date, err)
	}
	return &e, nil
}

func (s *pgEventStore) Create(ctx context.Context, name string, description string, start time.Time, end time.Time, volunteers int, location string) (*Event, error) {
	id := 0
	createdAt := time.Time{}

	query := "INSERT INTO events (name, description, start_date, end_date, volunteers_required, location) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id, created_at"
	err := s.db.QueryRowContext(ctx, query, name, description, start, end, volunteers, location).Scan(&id, &createdAt)
	if err != nil {
		return nil, fmt.Errorf("unable to create and store an event, error %w", err)
	}

	event := Event{id, name, description, start, end, volunteers, location, createdAt}
	return &event, nil
}

func (s *pgEventStore) Update(ctx context.Context, event Event) (*Event, error) {
	query := "UPDATE events SET name = $1, description = $2, start_date = $3, end_date = $4, volunteers_required = $5, location = $6, WHERE id = $7"

	_, err := s.db.ExecContext(ctx, query, event.Name, event.Description, event.StartDate, event.EndDate, event.VolunteersRequired, event.Location, event.ID)

	if err != nil {
		return nil, fmt.Errorf("unable to update event, error %w", err)

	}

	return &event, nil
}
func (s *pgEventStore) Delete(ctx context.Context, id int) error {
	query := "DELETE FROM events WHERE id = $1"

	_, err := s.db.ExecContext(ctx, query, id)

	if err != nil {
		return fmt.Errorf("unable to delete event, error %w", err)

	}

	return nil
}


// unregisters a user from an event 
func (s *pgEventStore) UnregisterEvent(ctx context.Context, userID int, eventID int) error {

	//returning id of the signUp from event_signups
	signUpQuery := `DELETE FROM event_signups WHERE user_id = $1 AND event_id = $2 RETURNING id`
	var signUpID int
	err := s.db.GetContext(ctx, &signUpID, signUpQuery, userID, eventID)

    if err != nil && err != sql.ErrNoRows {
        // Handle errors NOT corresponding to a 'no rows' return error 
        return fmt.Errorf("failed to remove user from event signup list: %w", err)
    }

	//returning id of the signUp from event_standby 
	if err == sql.ErrNoRows {
		var standByID int 
		standByQuery := `DELETE FROM event_standbys WHERE user_id = $1 AND event_id = $2 RETURNING id`
		err := s.db.GetContext(ctx, &standByID, standByQuery, userID, eventID)

		if err != nil {
			return fmt.Errorf("failed to remove user from event standby list: %w", err)
		}
	}
	return nil
}
