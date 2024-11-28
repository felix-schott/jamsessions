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
    genres VARCHAR(50)[] CHECK(genres <@ ARRAY['Straight-Ahead_Jazz'::VARCHAR, 'Modern_Jazz'::VARCHAR, 'Trad_Jazz'::VARCHAR, 'Jazz-Funk'::VARCHAR, 'Fusion'::VARCHAR, 'Latin_Jazz'::VARCHAR, 'Funk'::VARCHAR, 'RnB'::VARCHAR, 'Hip-Hop'::VARCHAR, 'Blues'::VARCHAR, 'Folk'::VARCHAR, 'Rock'::VARCHAR, 'Pop'::VARCHAR, 'World_Music'::VARCHAR]),
    start_time_utc TIMESTAMPTZ NOT NULL,
    interval VARCHAR(50) NOT NULL CHECK (interval IN ('Once', 'Daily', 'Weekly', 'Fortnightly', 'FirstOfMonth', 'SecondOfMonth', 'ThirdOfMonth', 'FourthOfMonth', 'LastOfMonth', 'IrregularWeekly')),
    duration_minutes SMALLINT NOT NULL,
    description TEXT NOT NULL,
	session_website VARCHAR(2000),
    dt_updated_utc TIMESTAMPTZ DEFAULT (NOW() AT TIME ZONE 'utc'),
    UNIQUE (venue, start_time_utc, interval) -- unique time and venue
);
-- create indices
CREATE INDEX jamsessions_venue_fkey_idx ON london_jam_sessions.jamsessions (venue);

-- TABLE london_jam_sessions.comments

CREATE TABLE london_jam_sessions.comments (
    comment_id SERIAL PRIMARY KEY,
    session INTEGER NOT NULL REFERENCES london_jam_sessions.jamsessions(session_id),
    author VARCHAR(200) NOT NULL,
    content TEXT NOT NULL,
    dt_posted TIMESTAMPTZ DEFAULT (NOW() AT TIME ZONE 'utc')
);
-- create indices
CREATE INDEX comments_session_fkey_idx ON london_jam_sessions.comments (session);

-- TABLE london_jam_sessions.ratings

CREATE TABLE london_jam_sessions.ratings (
    rating_id SERIAL PRIMARY KEY,
    session INTEGER NOT NULL REFERENCES london_jam_sessions.jamsessions(session_id),
    comment INTEGER REFERENCES london_jam_sessions.comments(comment_id), --optionally link to comment
    rating SMALLINT CHECK(rating < 6 AND rating > 0), -- 1 to 5
    dt_posted TIMESTAMPTZ DEFAULT (NOW() AT TIME ZONE 'utc')
);
-- create indices
CREATE INDEX ratings_session_fkey_idx ON london_jam_sessions.ratings (session);
CREATE INDEX ratings_comment_fkey_idx ON london_jam_sessions.ratings (comment);

-- create funcs to get session matches by date (range), used in queries
CREATE OR REPLACE FUNCTION london_jam_sessions.sessions_in_date_range(start date, stop date) 
RETURNS TABLE (session_id int, dates date[])
AS $$
    BEGIN
        RETURN QUERY
        SELECT s.session_id, array_agg(dates_in_range.d::date) AS dates
        FROM (
            SELECT generate_series(start, stop, interval '1 day') d
        ) dates_in_range
        CROSS JOIN london_jam_sessions.jamsessions s
        WHERE s.interval = 'Daily' AND s.start_time_utc::date <= stop
        OR (
            s.start_time_utc::date >= start
            AND start_time_utc::date < stop + interval '1 day'
            AND interval = 'Once'
        ) 
        OR (
            s.interval IN ('Weekly', 'IrregularWeekly') AND (dates_in_range.d::date - s.start_time_utc::date) % 7 = 0 
        ) 
        OR (
            s.interval = 'Fortnightly' AND (dates_in_range.d::date - s.start_time_utc::date) % 14 = 0 
        )
        OR (
            s.interval = 'FirstOfMonth' AND EXTRACT(dow FROM s.start_time_utc) = EXTRACT(dow FROM dates_in_range.d::date)
            AND CEIL (EXTRACT(DAY FROM dates_in_range.d) / 7) = 1
        )
        OR (
            s.interval = 'SecondOfMonth' AND EXTRACT(dow FROM s.start_time_utc) = EXTRACT(dow FROM dates_in_range.d::date)
            AND CEIL (EXTRACT(DAY FROM dates_in_range.d) / 7) = 2
        )
        OR (
            s.interval = 'ThirdOfMonth' AND EXTRACT(dow FROM s.start_time_utc) = EXTRACT(dow FROM dates_in_range.d::date)
            AND CEIL (EXTRACT(DAY FROM dates_in_range.d) / 7) = 3
        )
        OR (
            s.interval = 'FourthOfMonth' AND EXTRACT(dow FROM s.start_time_utc) = EXTRACT(dow FROM dates_in_range.d::date)
            AND CEIL (EXTRACT(DAY FROM dates_in_range.d) / 7) = 4
        )
        OR s.interval = 'LastOfMonth' AND EXTRACT(dow FROM s.start_time_utc) = EXTRACT(dow FROM dates_in_range.d::date)
        AND EXTRACT(MONTH FROM (dates_in_range.d + interval '7 day')) != EXTRACT(MONTH FROM dates_in_range.d)
        GROUP BY s.session_id;
    END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION london_jam_sessions.sessions_on_date(d date) 
RETURNS TABLE (session_id int, dates date[])
AS $$
    BEGIN
        RETURN QUERY
        SELECT s.session_id, ARRAY[d] FROM london_jam_sessions.jamsessions s
        WHERE s.start_time_utc::date <= d -- make sure we don't match any future sessions
        AND (
            s.interval = 'Daily' AND s.start_time_utc::date <= d
        )
        OR (
            s.start_time_utc::date = d::date
            AND s.interval = 'Once'
        ) 
        OR (
            s.interval IN ('Weekly', 'IrregularWeekly') AND EXTRACT(dow FROM s.start_time_utc) = EXTRACT(dow FROM d::date)
        ) 
        OR (
            s.interval = 'Fortnightly' AND (d::date - s.start_time_utc::date) % 14 = 0 -- check if the number of days between is divisible by 14
        )
        OR (
            s.interval = 'FirstOfMonth' AND EXTRACT(dow FROM s.start_time_utc) = EXTRACT(dow FROM d::date) 
            AND CEIL (
                EXTRACT(DAY FROM d) / 7
            ) = 1
        )
        OR (
            s.interval = 'SecondOfMonth' AND EXTRACT(dow FROM s.start_time_utc) = EXTRACT(dow FROM d::date)
            AND CEIL (
                EXTRACT(DAY FROM d) / 7
            ) = 2
        )
        OR (
            s.interval = 'ThirdOfMonth' AND EXTRACT(dow FROM s.start_time_utc) = EXTRACT(dow FROM d::date) 
            AND CEIL (
                EXTRACT(DAY FROM d) / 7
            ) = 3
        )
        OR (
            s.interval = 'FourthOfMonth' AND EXTRACT(dow FROM s.start_time_utc) = EXTRACT(dow FROM d::date) 
            AND CEIL (
                EXTRACT(DAY FROM d) / 7
            ) = 4
        )
        OR (
            s.interval = 'LastOfMonth' AND EXTRACT(dow FROM s.start_time_utc) = EXTRACT(dow FROM d::date) 
            AND (
                EXTRACT(MONTH FROM (d::date + interval '7 day')) != EXTRACT(MONTH FROM d) 
            )
        );
    END;
$$ LANGUAGE plpgsql;