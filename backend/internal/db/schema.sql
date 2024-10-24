CREATE EXTENSION postgis;

-- create schema

CREATE SCHEMA london_jam_sessions AUTHORIZATION postgres;

-- create london_jam_sessions.venues table

CREATE TABLE london_jam_sessions.venues (
    venue_id SERIAL PRIMARY KEY,
    venue_name VARCHAR(250) NOT NULL UNIQUE, -- name must be unique
    address_first_line VARCHAR(100) NOT NULL,
    address_second_line VARCHAR(100),
    city VARCHAR(200) NOT NULL,
    postcode VARCHAR(8) NOT NULL,
    geom GEOMETRY(Point, 4326) NOT NULL,
	venue_website VARCHAR(2000),
    backline VARCHAR(20)[] CHECK(backline <@ ARRAY['PA'::VARCHAR, 'Guitar_Amp'::VARCHAR, 'Bass_Amp'::VARCHAR, 'Keys'::VARCHAR, 'Drums'::VARCHAR, 'Microphone'::VARCHAR, 'MiscPercussion'::VARCHAR]),
    venue_comments TEXT[],
    venue_dt_updated_utc TIMESTAMPTZ DEFAULT (NOW() AT TIME ZONE 'utc'),
    UNIQUE (address_first_line, postcode) -- unique address, there can't be two london_jam_sessions.venues at the same address
);
-- create indices
CREATE INDEX venues_venue_name_idx ON london_jam_sessions.venues (venue_name);
CREATE INDEX venues_backline_idx ON london_jam_sessions.venues USING GIN (backline);

-- trigger to propagate dt_updated to london_jam_sessions.jamsessions table
-- every time the london_jam_sessions.venues table is updated, the timestamp of the corresponding sessions is updated too
CREATE FUNCTION update_timestamp_venue() RETURNS trigger AS $$
    BEGIN
        UPDATE london_jam_sessions.jamsessions SET dt_updated_utc = NEW.venue_dt_updated_utc WHERE venue = NEW.venue_id;
        RETURN NEW;
    END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER update_timestamp_venue AFTER INSERT OR UPDATE ON london_jam_sessions.venues
    FOR EACH ROW EXECUTE FUNCTION update_timestamp_venue();

-- TABLE london_jam_sessions.jamsessions

CREATE TABLE london_jam_sessions.jamsessions (
    session_id SERIAL PRIMARY KEY,
    session_name VARCHAR(250) NOT NULL,
    venue INTEGER NOT NULL REFERENCES london_jam_sessions.venues(venue_id) ON DELETE CASCADE, -- if a venue is deleted, all sessions associated with the venue should be deleted too
    genres VARCHAR(50)[] CHECK(genres <@ ARRAY['Straight-Ahead_Jazz'::VARCHAR, 'Modern_Jazz'::VARCHAR, 'Trad_Jazz'::VARCHAR, 'Jazz-Funk'::VARCHAR, 'Fusion'::VARCHAR, 'Latin_Jazz'::VARCHAR, 'Funk'::VARCHAR, 'Blues'::VARCHAR, 'Folk'::VARCHAR, 'Rock'::VARCHAR, 'World_Music'::VARCHAR]),
    start_time_utc TIMESTAMPTZ NOT NULL,
    interval VARCHAR(13) NOT NULL CHECK (interval IN ('Once', 'Daily', 'Weekly', 'FirstOfMonth', 'SecondOfMonth', 'ThirdOfMonth', 'FourthOfMonth', 'LastOfMonth')),
    duration_minutes SMALLINT NOT NULL,
    description TEXT NOT NULL,
    session_comments TEXT[],
	session_website VARCHAR(2000),
    dt_updated_utc TIMESTAMPTZ DEFAULT (NOW() AT TIME ZONE 'utc'),
    UNIQUE (venue, start_time_utc, interval) -- unique time and venue
);
-- create indices
CREATE INDEX jamsessions_venue_fkey_idx ON london_jam_sessions.jamsessions (venue);