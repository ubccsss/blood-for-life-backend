--- Events
CREATE TABLE IF NOT EXISTS events (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    start_date TIMESTAMP NOT NULL,
    end_date TIMESTAMP NOT NULL,
    volunteers_required INTEGER NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
);

--- Users
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    UNIQUE(email)
);

--- Event Attendance
--- 
CREATE TABLE IF NOT EXISTS event_attendance (
    id SERIAL PRIMARY KEY,
    event_id INTEGER NOT NULL REFERENCES events(id),
    user_id INTEGER NOT NULL REFERENCES users(id),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    UNIQUE(event_id, user_id)
);

--- Event Signups
--- Store information about who has signed up for an event
--- There should be no more than volunteers_required rows for a given event
CREATE TABLE IF NOT EXISTS event_signups (
    id SERIAL PRIMARY KEY,
    event_id INTEGER NOT NULL REFERENCES events(id),
    user_id INTEGER NOT NULL REFERENCES users(id),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    UNIQUE(event_id, user_id)
);

--- Event Standbys
--- Store information about who is on standby for an event
--- There can in theory be an unlimited number of standbys for a given event
CREATE TABLE IF NOT EXISTS event_standbys (
    id SERIAL PRIMARY KEY,
    event_id INTEGER NOT NULL REFERENCES events(id),
    user_id INTEGER NOT NULL REFERENCES users(id),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    UNIQUE(event_id, user_id)
);