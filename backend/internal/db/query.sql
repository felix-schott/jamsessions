-- name: GetAllVenuesAsGeoJSON :one
SELECT json_build_object(
    'type', 'FeatureCollection',
    'features', json_agg(public.ST_AsGeoJSON(v.*)::json)
) FROM london_jam_sessions.venues v;

-- name: GetVenueById :one
SELECT * FROM london_jam_sessions.venues
WHERE venue_id = $1;

-- name: GetVenueByIdAsGeoJSON :one
WITH t AS (
    SELECT * FROM london_jam_sessions.venues
    WHERE venue_id = $1
)
SELECT public.ST_AsGeoJSON(t.*) FROM t;

-- name: GetVenueByName :one
SELECT * FROM london_jam_sessions.venues
WHERE venue_name = $1;

-- name: GetAllSessions :many
SELECT * FROM london_jam_sessions.jamsessions s
JOIN london_jam_sessions.venues l ON s.venue = l.venue_id;

-- name: GetAllSessionsAsGeoJSON :one
WITH t AS (
    SELECT * FROM london_jam_sessions.jamsessions s
    JOIN london_jam_sessions.venues l ON s.venue = l.venue_id
)
SELECT json_build_object(
    'type', 'FeatureCollection',
    'features', json_agg(public.ST_AsGeoJSON(t.*)::json)
) FROM t;

-- name: GetSessionById :one
SELECT * FROM london_jam_sessions.jamsessions s
JOIN london_jam_sessions.venues l ON s.venue = l.venue_id
WHERE s.session_id = $1;

-- name: GetSessionByIdAsGeoJSON :one
WITH t AS (
    SELECT * FROM london_jam_sessions.jamsessions s
    JOIN london_jam_sessions.venues l ON s.venue = l.venue_id
    WHERE s.session_id = $1
)
SELECT public.ST_AsGeoJSON(t.*) FROM t;

-- name: GetSessionsByDateAsGeoJSON :one
WITH t AS (
    SELECT * FROM london_jam_sessions.jamsessions s
    JOIN london_jam_sessions.venues l ON s.venue = l.venue_id
    WHERE interval = 'Daily' 
    OR (
        start_time_utc::date = sqlc.arg(date)::date
        AND interval = 'Once'
    ) 
    OR (
        interval = 'Weekly' AND EXTRACT(dow FROM s.start_time_utc) = EXTRACT(dow FROM sqlc.arg(date)::date)
    ) 
    OR (
        interval = 'Fortnightly' AND EXTRACT(DAY FROM sqlc.arg(date) - start_time_utc) % 14 = 0
    )
    OR (
        interval = 'FirstOfMonth' AND CEIL (
            EXTRACT(DAY FROM sqlc.arg(date)) / 7
        ) = 1
    )
    OR (
        interval = 'SecondOfMonth' AND CEIL (
            EXTRACT(DAY FROM sqlc.arg(date)) / 7
        ) = 2
    )
    OR (
        interval = 'ThirdOfMonth' AND CEIL (
            EXTRACT(DAY FROM sqlc.arg(date)) / 7
        ) = 3
    )
    OR (
        interval = 'FourthOfMonth' AND CEIL (
            EXTRACT(DAY FROM sqlc.arg(date)) / 7
        ) = 4
    )
    OR interval = 'LastOfMonth' AND (
        EXTRACT(MONTH FROM sqlc.arg(date) + 7) != EXTRACT(MONTH FROM sqlc.arg(date)) 
    )
)
SELECT json_build_object(
    'type', 'FeatureCollection',
    'features', json_agg(public.ST_AsGeoJSON(t.*)::json)
) FROM t;

-- name: GetSessionsByBacklineAsGeoJSON :one
WITH t AS (
    SELECT * FROM london_jam_sessions.jamsessions s
    JOIN london_jam_sessions.venues l ON s.venue = l.venue_id
    WHERE l.backline @> $1
)
SELECT json_build_object(
    'type', 'FeatureCollection',
    'features', json_agg(public.ST_AsGeoJSON(t.*)::json)
) FROM t;

-- name: GetSessionsByGenreAsGeoJSON :one
WITH t AS (
    SELECT * FROM london_jam_sessions.jamsessions s
    JOIN london_jam_sessions.venues l ON s.venue = l.venue_id
    WHERE s.genres @> $1
)
SELECT json_build_object(
    'type', 'FeatureCollection',
    'features', json_agg(public.ST_AsGeoJSON(*)::json)
) FROM t;

-- name: GetSessionsByDateAndGenreAsGeoJSON :one
WITH t AS (
    SELECT * FROM london_jam_sessions.jamsessions s
    JOIN london_jam_sessions.venues l ON s.venue = l.venue_id
    WHERE s.genres @> sqlc.arg(genres) 
    AND interval = 'Daily'
    OR (
        start_time_utc::date = sqlc.arg(date)::date
        AND interval = 'Once'
    ) 
    OR (
        interval = 'Weekly' AND EXTRACT(dow FROM s.start_time_utc) = EXTRACT(dow FROM sqlc.arg(date)::date)
    ) 
    OR (
        interval = 'Fortnightly' AND EXTRACT(DAY FROM sqlc.arg(date) - start_time_utc) % 14 = 0
    )
    OR (
        interval = 'FirstOfMonth' AND CEIL (
            EXTRACT(DAY FROM sqlc.arg(date)) / 7
        ) = 1
    )
    OR (
        interval = 'SecondOfMonth' AND CEIL (
            EXTRACT(DAY FROM sqlc.arg(date)) / 7
        ) = 2
    )
    OR (
        interval = 'ThirdOfMonth' AND CEIL (
            EXTRACT(DAY FROM sqlc.arg(date)) / 7
        ) = 3
    )
    OR (
        interval = 'FourthOfMonth' AND CEIL (
            EXTRACT(DAY FROM sqlc.arg(date)) / 7
        ) = 4
    )
    OR interval = 'LastOfMonth' AND (
        EXTRACT(MONTH FROM sqlc.arg(date) + 7) != EXTRACT(MONTH FROM sqlc.arg(date)) 
    )
)
SELECT json_build_object(
    'type', 'FeatureCollection',
    'features', json_agg(public.ST_AsGeoJSON(t.*)::json)
) FROM t;

-- name: GetSessionsByDateAndBacklineAsGeoJSON :one
WITH t AS (
    SELECT * FROM london_jam_sessions.jamsessions s
    JOIN london_jam_sessions.venues l ON s.venue = l.venue_id
    WHERE l.backline @> sqlc.arg(backline)
    AND interval = 'Daily'
    OR (
        start_time_utc::date = sqlc.arg(date)::date
        AND interval = 'Once'
    ) 
    OR (
        interval = 'Weekly' AND EXTRACT(dow FROM s.start_time_utc) = EXTRACT(dow FROM sqlc.arg(date)::date)
    ) 
    OR (
        interval = 'Fortnightly' AND EXTRACT(DAY FROM sqlc.arg(date) - start_time_utc) % 14 = 0
    )
    OR (
        interval = 'FirstOfMonth' AND CEIL (
            EXTRACT(DAY FROM sqlc.arg(date)) / 7
        ) = 1
    )
    OR (
        interval = 'SecondOfMonth' AND CEIL (
            EXTRACT(DAY FROM sqlc.arg(date)) / 7
        ) = 2
    )
    OR (
        interval = 'ThirdOfMonth' AND CEIL (
            EXTRACT(DAY FROM sqlc.arg(date)) / 7
        ) = 3
    )
    OR (
        interval = 'FourthOfMonth' AND CEIL (
            EXTRACT(DAY FROM sqlc.arg(date)) / 7
        ) = 4
    )
    OR interval = 'LastOfMonth' AND (
        EXTRACT(MONTH FROM sqlc.arg(date) + 7) != EXTRACT(MONTH FROM sqlc.arg(date)) 
    )
)
SELECT json_build_object(
    'type', 'FeatureCollection',
    'features', json_agg(public.ST_AsGeoJSON(t.*)::json)
) FROM t;

-- name: GetSessionsByDateAndGenreAndBacklineAsGeoJSON :one
WITH t AS (
    SELECT * FROM london_jam_sessions.jamsessions s
    JOIN london_jam_sessions.venues l ON s.venue = l.venue_id
    WHERE s.genres @> sqlc.arg(genres)
    AND l.backline @> sqlc.arg(backline)
    AND interval = 'Daily'
    OR (
        start_time_utc::date = sqlc.arg(date)::date
        AND interval = 'Once'
    ) 
    OR (
        interval = 'Weekly' AND EXTRACT(dow FROM s.start_time_utc) = EXTRACT(dow FROM sqlc.arg(date)::date)
    ) 
    OR (
        interval = 'Fortnightly' AND EXTRACT(DAY FROM sqlc.arg(date) - start_time_utc) % 14 = 0
    )
    OR (
        interval = 'FirstOfMonth' AND CEIL (
            EXTRACT(DAY FROM sqlc.arg(date)) / 7
        ) = 1
    )
    OR (
        interval = 'SecondOfMonth' AND CEIL (
            EXTRACT(DAY FROM sqlc.arg(date)) / 7
        ) = 2
    )
    OR (
        interval = 'ThirdOfMonth' AND CEIL (
            EXTRACT(DAY FROM sqlc.arg(date)) / 7
        ) = 3
    )
    OR (
        interval = 'FourthOfMonth' AND CEIL (
            EXTRACT(DAY FROM sqlc.arg(date)) / 7
        ) = 4
    )
    OR interval = 'LastOfMonth' AND (
        EXTRACT(MONTH FROM sqlc.arg(date) + 7) != EXTRACT(MONTH FROM sqlc.arg(date)) 
    )
)
SELECT json_build_object(
    'type', 'FeatureCollection',
    'features', json_agg(public.ST_AsGeoJSON(t.*)::json)
) FROM t;

-- name: GetSessionsByGenreAndBacklineAsGeoJSON :one
WITH t AS (
    SELECT * FROM london_jam_sessions.jamsessions s
    JOIN london_jam_sessions.venues l ON s.venue = l.venue_id
    WHERE s.genres @> sqlc.arg(genres)
    AND l.backline @> sqlc.arg(backline)
)
SELECT json_build_object(
    'type', 'FeatureCollection',
    'features', json_agg(public.ST_AsGeoJSON(t.*)::json)
) FROM t;

-- name: InsertVenue :one
INSERT INTO london_jam_sessions.venues (
    venue_name, address_first_line, address_second_line, city, postcode, geom, venue_website, backline, venue_comments
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9
) RETURNING venue_id;

-- name: UpdateVenueById :exec
UPDATE london_jam_sessions.venues
SET -- see https://docs.sqlc.dev/en/latest/howto/named_parameters.html#nullable-parameters
    venue_name = coalesce(sqlc.narg(venue_name), venue_name),
    address_first_line = coalesce(sqlc.narg(address_first_line), address_first_line),
    address_second_line = coalesce(sqlc.narg(address_second_line), address_second_line),
    city = coalesce(sqlc.narg(city), city),
    postcode = coalesce(sqlc.narg(postcode), postcode),
    geom = coalesce(sqlc.narg(geom), geom),
    venue_website = coalesce(sqlc.narg(venue_website), venue_website),
    backline = coalesce(sqlc.narg(backline), backline),
    venue_comments = coalesce(sqlc.narg(venue_comments), venue_comments)
WHERE venue_id = $1;
-- SET venue_name = CASE WHEN sqlc.arg(set_venue_name)::bool
--         THEN sqlc.arg(venue_name)
--         ELSE venue_name
--         END,
--     address_first_line = CASE WHEN sqlc.arg(set_address_first_line)::bool
--         THEN sqlc.arg(address_first_line)
--         ELSE address_first_line
--         END, 
--     address_second_line = CASE WHEN sqlc.arg(set_address_second_line)::bool
--         THEN sqlc.arg(address_second_line)
--         ELSE address_second_line
--         END, 
--     city = CASE WHEN sqlc.arg(set_city)::bool
--         THEN sqlc.arg(city)
--         ELSE city
--         END, 
--     postcode = CASE WHEN sqlc.arg(set_postcode)::bool
--         THEN sqlc.arg(postcode)
--         ELSE postcode
--         END,
--     geom = CASE WHEN sqlc.arg(set_geom)::bool
--         THEN sqlc.arg(geom)
--         ELSE geom
--         END, 
--     venue_website = CASE WHEN sqlc.arg(set_venue_website)::bool
--         THEN sqlc.arg(venue_website)
--         ELSE venue_website
--         END,  
--     backline = CASE WHEN sqlc.arg(set_backline)::bool
--         THEN sqlc.arg(backline)
--         ELSE backline
--         END,  
--     venue_comments = CASE WHEN sqlc.arg(set_venue_comments)::bool
--         THEN sqlc.arg(venue_comments)
--         ELSE venue_comments
--         END
-- WHERE venue_id = $1;

-- name: InsertJamSession :one
INSERT INTO london_jam_sessions.jamsessions (
    session_name, venue, description, genres, start_time_utc, interval, duration_minutes, session_website, session_comments
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9
) RETURNING session_id;

-- name: UpdateJamSessionById :exec
UPDATE london_jam_sessions.jamsessions
SET
    session_name = coalesce(sqlc.narg(session_name), session_name),
    description = coalesce(sqlc.narg(description), description),
    genres = coalesce(sqlc.narg(genres), genres),
    start_time_utc = coalesce(sqlc.narg(start_time_utc), start_time_utc),
    interval = coalesce(sqlc.narg(interval), interval),
    duration_minutes = coalesce(sqlc.narg(duration_minutes), duration_minutes),
    session_website = coalesce(sqlc.narg(session_website), session_website),
    session_comments = coalesce(sqlc.narg(session_comments), session_comments)
WHERE session_id = $1;

-- SET session_name = CASE WHEN sqlc.arg(set_session_name)::bool
--         THEN sqlc.arg(session_name)
--         ELSE session_name
--         END,
--     description = CASE WHEN sqlc.arg(set_description)::bool
--         THEN sqlc.arg(description)
--         ELSE description
--         END,
--     genres = CASE WHEN sqlc.arg(set_genres)::bool
--         THEN sqlc.arg(genres)
--         ELSE genres
--         END,
--     start_time_utc = CASE WHEN sqlc.arg(set_start_time_utc)::bool
--         THEN sqlc.arg(start_time_utc)
--         ELSE start_time_utc
--         END,
--     interval = CASE WHEN sqlc.arg(set_interval)::bool
--         THEN sqlc.arg(interval)
--         ELSE interval
--         END,
--     duration_minutes = CASE WHEN sqlc.arg(set_duration_minutes)::bool
--         THEN sqlc.arg(duration_minutes)
--         ELSE duration_minutes
--         END, 
--     session_website = CASE WHEN sqlc.arg(set_session_website)::bool
--         THEN sqlc.arg(session_website)
--         ELSE session_website
--         END,
--     session_comments = CASE WHEN sqlc.arg(set_session_comments)::bool
--         THEN sqlc.arg(session_comments)
--         ELSE session_comments
--         END 
-- WHERE session_id = $1;

-- name: DeleteJamSessionById :exec
DELETE FROM london_jam_sessions.jamsessions
WHERE session_id = $1;

-- name: DeleteVenueById :exec
DELETE FROM london_jam_sessions.venues
WHERE venue_id = $1;

-- name: DeleteVenueByJamSessionId :exec
DELETE FROM london_jam_sessions.venues l
USING london_jam_sessions.jamsessions s
WHERE s.venue = l.venue_id AND s.session_id = $1;