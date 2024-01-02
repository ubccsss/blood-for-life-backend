package store
import (
    "context"
    "fmt"
    "time"
    "github.com/jmoiron/sqlx"
)

type Event struct {
    ID                  int       `db:"id" json:"id"`
    Name                string    `db:"name" json:"name"`
    Description         string    `db:"description" json:"description"`
    StartDate           time.Time `db:"start_date" json:"startDate"`
    EndDate             time.Time `db:"end_date" json:"endDate"`
    VolunteersRequired  int       `db:"volunteers_required" json:"volunteersRequired"`
    CreatedAt           time.Time `db:"created_at" json:"createdAt"`
}

type EventStore interface {
    GetAll(ctx context.Context) ([]Event, error)
    GetOne(ctx context.Context, id int) (*Event, error)
    GetOneByStartDate(ctx context.Context, startDate time.Time) (*Event, error)
    GetOneByName(ctx context.Context, name string) (*Event, error) // ??
    Create(ctx context.Context, event Event) (*Event, error)
    Update(ctx context.Context, event Event) (*Event, error)
    Delete(ctx context.Context, id int) error
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
    err := s.db.GetContext(ctx, &e, "SELECT * FROM events WHERE start_date = $1", date)
    if err != nil {
        return nil, fmt.Errorf("unable to find event with start date %s, with error %w", date, err)
    }
    return &e, nil
}

func (s *pgEventStore) Create(ctx context.Context, event Event) (*Event, error) {
    query := "INSERT INTO events (name, description, start_date, end_date, volunteers_required) VALUES ($1, $2, $3, $4, $5) RETURNING id, created_at"
    err := s.db.QueryRowContext(ctx, query, event.Name, event.Description, event.StartDate, event.EndDate, event.VolunteersRequired).Scan(&event.ID, &event.CreatedAt)
    if err != nil {
        return nil, fmt.Errorf("unable to create and store an event, error %w", err)
    }
    return &event, nil 
}

func (s *pgEventStore) Update(ctx context.Context, event Event) (*Event, error) {
	query := "UPDATE events SET name = $1, description = $2, start_date = $3, end_date = $4, volunteers_required = $5, WHERE id = $6"

	_, err = s.db.ExecContext(ctx, query, event.Name, event.Description, event.StartDate, event.EndDate, event.VolunteersRequired, event.ID)

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