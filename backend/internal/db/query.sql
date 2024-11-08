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

-- name: GetCommentsBySessionId :many
SELECT c.*, r.rating FROM london_jam_sessions.comments c
LEFT OUTER JOIN london_jam_sessions.ratings r ON c.comment_id = r.rating_id
WHERE c.session = $1;

-- name: GetAllSessions :many
SELECT s.*, l.*, coalesce(round(avg(rating), 2), 0)::real AS rating FROM london_jam_sessions.jamsessions s
JOIN london_jam_sessions.venues l ON s.venue = l.venue_id
LEFT OUTER JOIN london_jam_sessions.ratings r ON s.session_id = r.session
GROUP BY s.session_id, l.venue_id;

-- name: GetAllSessionsAsGeoJSON :one
WITH t AS (
    SELECT s.*, l.*, coalesce(round(avg(rating), 2), 0)::real AS rating FROM london_jam_sessions.jamsessions s
    JOIN london_jam_sessions.venues l ON s.venue = l.venue_id
    LEFT OUTER JOIN london_jam_sessions.ratings r ON s.session_id = r.session
    GROUP BY s.session_id, l.venue_id
)
SELECT json_build_object(
    'type', 'FeatureCollection',
    'features', json_agg(public.ST_AsGeoJSON(t.*)::json)
) FROM t;

-- name: GetSessionById :one
SELECT s.*, l.*, coalesce(round(avg(rating), 2), 0.0)::real AS rating FROM london_jam_sessions.jamsessions s
JOIN london_jam_sessions.venues l ON s.venue = l.venue_id
LEFT OUTER JOIN london_jam_sessions.ratings r ON s.session_id = r.session
WHERE s.session_id = $1
GROUP BY s.session_id, l.venue_id;

-- name: GetSessionByIdAsGeoJSON :one
WITH t AS (
    SELECT s.*, l.*, coalesce(round(avg(rating), 2), 0.0)::real AS rating FROM london_jam_sessions.jamsessions s
    JOIN london_jam_sessions.venues l ON s.venue = l.venue_id
    LEFT OUTER JOIN london_jam_sessions.ratings r ON s.session_id = r.session
    WHERE s.session_id = $1
    GROUP BY s.session_id, l.venue_id
)
SELECT public.ST_AsGeoJSON(t.*) FROM t;

-- name: GetSessionsByDateAsGeoJSON :one
WITH t AS (
    SELECT s.*, l.*, coalesce(round(avg(rating), 2), 0.0)::real AS rating 
    FROM london_jam_sessions.sessions_on_date(sqlc.arg(date)::date) d
    LEFT OUTER JOIN london_jam_sessions.jamsessions s ON d.session_id = s.session_id
    JOIN london_jam_sessions.venues l ON s.venue = l.venue_id
    LEFT OUTER JOIN london_jam_sessions.ratings r ON s.session_id = r.session
    GROUP BY s.session_id, l.venue_id
)
SELECT json_build_object(
    'type', 'FeatureCollection',
    'features', json_agg(public.ST_AsGeoJSON(t.*)::json)
) FROM t;

-- name: GetSessionIdsByDateRange :many
SELECT * FROM london_jam_sessions.sessions_in_date_range(sqlc.arg(start_date)::date, sqlc.arg(end_date)::date);

-- name: GetSessionIdsByDate :many
SELECT * FROM london_jam_sessions.sessions_on_date(sqlc.arg(date)::date);

-- name: GetSessionsByDateRangeAsGeoJSON :one
WITH t AS (
    SELECT d.d, s.*, l.*, coalesce(round(avg(rating), 2), 0.0)::real AS rating 
    FROM london_jam_sessions.jamsessions s 
    LEFT OUTER JOIN london_jam_sessions.sessions_in_date_range(sqlc.arg(start_date)::date, sqlc.arg(end_date)::date) d 
    ON d.session_id = s.session_id 
    JOIN london_jam_sessions.venues l ON s.venue = l.venue_id
    LEFT OUTER JOIN london_jam_sessions.ratings r ON s.session_id = r.session
    GROUP BY s.session_id, l.venue_id, d.d
)
SELECT json_build_object(
    'type', 'FeatureCollection',
    'features', json_agg(public.ST_AsGeoJSON(t.*)::json)
) FROM t;

-- name: GetSessionsByBacklineAsGeoJSON :one
WITH t AS (
    SELECT s.*, l.*, coalesce(round(avg(rating), 2), 0.0)::real AS rating FROM london_jam_sessions.jamsessions s
    JOIN london_jam_sessions.venues l ON s.venue = l.venue_id
    LEFT OUTER JOIN london_jam_sessions.ratings r ON s.session_id = r.session
    WHERE l.backline @> $1
    GROUP BY s.session_id, l.venue_id
)
SELECT json_build_object(
    'type', 'FeatureCollection',
    'features', json_agg(public.ST_AsGeoJSON(t.*)::json)
) FROM t;

-- name: GetSessionsByGenreAsGeoJSON :one
WITH t AS (
    SELECT s.*, l.*, coalesce(round(avg(rating), 2), 0.0)::real AS rating FROM london_jam_sessions.jamsessions s
    JOIN london_jam_sessions.venues l ON s.venue = l.venue_id
    LEFT OUTER JOIN london_jam_sessions.ratings r ON s.session_id = r.session
    WHERE s.genres @> $1
    GROUP BY s.session_id, l.venue_id
)
SELECT json_build_object(
    'type', 'FeatureCollection',
    'features', json_agg(public.ST_AsGeoJSON(*)::json)
) FROM t;

-- name: GetSessionsByDateAndGenreAsGeoJSON :one
WITH t AS (
    SELECT s.*, l.*, coalesce(round(avg(rating), 2), 0.0)::real AS rating FROM london_jam_sessions.jamsessions s
    LEFT OUTER JOIN london_jam_sessions.sessions_on_date(sqlc.arg(date)::date) d ON d.session_id = s.session_id 
    JOIN london_jam_sessions.venues l ON s.venue = l.venue_id
    LEFT OUTER JOIN london_jam_sessions.ratings r ON s.session_id = r.session
    WHERE s.genres @> sqlc.arg(genres) 
    GROUP BY s.session_id, l.venue_id
)
SELECT json_build_object(
    'type', 'FeatureCollection',
    'features', json_agg(public.ST_AsGeoJSON(t.*)::json)
) FROM t;

-- name: GetSessionsByDateRangeAndGenreAsGeoJSON :one
WITH t AS (
    SELECT d.d, s.*, l.*, coalesce(round(avg(rating), 2), 0.0)::real AS rating 
    FROM london_jam_sessions.jamsessions s 
    LEFT OUTER JOIN london_jam_sessions.sessions_in_date_range(sqlc.arg(start_date)::date, sqlc.arg(end_date)::date) d 
    ON d.session_id = s.session_id 
    JOIN london_jam_sessions.venues l ON s.venue = l.venue_id
    LEFT OUTER JOIN london_jam_sessions.ratings r ON s.session_id = r.session
    WHERE s.genres @> sqlc.arg(genres)
    GROUP BY s.session_id, l.venue_id, d.d
)
SELECT json_build_object(
    'type', 'FeatureCollection',
    'features', json_agg(public.ST_AsGeoJSON(t.*)::json)
) FROM t;


-- name: GetSessionsByDateAndBacklineAsGeoJSON :one
WITH t AS (
    SELECT s.*, l.*, coalesce(round(avg(rating), 2), 0.0)::real AS rating FROM london_jam_sessions.jamsessions s
    LEFT OUTER JOIN london_jam_sessions.sessions_on_date(sqlc.arg(date)::date) d ON d.session_id = s.session_id
    JOIN london_jam_sessions.venues l ON s.venue = l.venue_id
    LEFT OUTER JOIN london_jam_sessions.ratings r ON s.session_id = r.session
    WHERE l.backline @> sqlc.arg(backline)
    GROUP BY s.session_id, l.venue_id
)
SELECT json_build_object(
    'type', 'FeatureCollection',
    'features', json_agg(public.ST_AsGeoJSON(t.*)::json)
) FROM t;

-- name: GetSessionsByDateRangeAndBacklineAsGeoJSON :one
WITH t AS (
    SELECT d.d, s.*, l.*, coalesce(round(avg(rating), 2), 0.0)::real AS rating 
    FROM london_jam_sessions.jamsessions s 
    LEFT OUTER JOIN london_jam_sessions.sessions_in_date_range(sqlc.arg(start_date)::date, sqlc.arg(end_date)::date) d 
    ON d.session_id = s.session_id 
    JOIN london_jam_sessions.venues l ON s.venue = l.venue_id
    LEFT OUTER JOIN london_jam_sessions.ratings r ON s.session_id = r.session
    WHERE l.backline @> sqlc.arg(backline)
    GROUP BY s.session_id, l.venue_id, d.d
)
SELECT json_build_object(
    'type', 'FeatureCollection',
    'features', json_agg(public.ST_AsGeoJSON(t.*)::json)
) FROM t;

-- name: GetSessionsByDateAndGenreAndBacklineAsGeoJSON :one
WITH t AS (
    SELECT s.*, l.*, coalesce(round(avg(rating), 2), 0.0)::real AS rating FROM london_jam_sessions.jamsessions s
    LEFT OUTER JOIN london_jam_sessions.sessions_on_date(sqlc.arg(date)::date) d ON d.session_id = s.session_id
    JOIN london_jam_sessions.venues l ON s.venue = l.venue_id
    LEFT OUTER JOIN london_jam_sessions.ratings r ON s.session_id = r.session
    WHERE s.genres @> sqlc.arg(genres)
    AND l.backline @> sqlc.arg(backline)
    GROUP BY s.session_id,l.venue_id
)
SELECT json_build_object(
    'type', 'FeatureCollection',
    'features', json_agg(public.ST_AsGeoJSON(t.*)::json)
) FROM t;

-- name: GetSessionsByDateRangeAndGenreAndBacklineAsGeoJSON :one
WITH t AS (
    SELECT d.d, s.*, l.*, coalesce(round(avg(rating), 2), 0.0)::real AS rating 
    FROM london_jam_sessions.jamsessions s 
    LEFT OUTER JOIN london_jam_sessions.sessions_in_date_range(sqlc.arg(start_date)::date, sqlc.arg(end_date)::date) d 
    ON d.session_id = s.session_id 
    JOIN london_jam_sessions.venues l ON s.venue = l.venue_id
    LEFT OUTER JOIN london_jam_sessions.ratings r ON s.session_id = r.session
    WHERE s.genres @> sqlc.arg(genres)
    AND l.backline @> sqlc.arg(backline)
    GROUP BY s.session_id, l.venue_id, d.d
)
SELECT json_build_object(
    'type', 'FeatureCollection',
    'features', json_agg(public.ST_AsGeoJSON(t.*)::json)
) FROM t;

-- name: GetSessionsByGenreAndBacklineAsGeoJSON :one
WITH t AS (
    SELECT s.*, l.*, coalesce(round(avg(rating), 2), 0.0)::real AS rating FROM london_jam_sessions.jamsessions s
    JOIN london_jam_sessions.venues l ON s.venue = l.venue_id
    LEFT OUTER JOIN london_jam_sessions.ratings r ON s.session_id = r.session
    WHERE s.genres @> sqlc.arg(genres)
    AND l.backline @> sqlc.arg(backline)
    GROUP BY s.session_id, l.venue_id
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

-- name: InsertJamSession :one
INSERT INTO london_jam_sessions.jamsessions (
    session_name, venue, description, genres, start_time_utc, interval, duration_minutes, session_website
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8
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
    session_website = coalesce(sqlc.narg(session_website), session_website)
WHERE session_id = $1;

-- name: InsertSessionComment :one
INSERT INTO london_jam_sessions.comments (
    session, author, content
) VALUES (
    $1, $2, $3
) RETURNING comment_id;

-- name: InsertSessionRating :one
INSERT INTO london_jam_sessions.ratings (
    session, rating, comment
) VALUES (
    $1, $2, $3
) RETURNING rating_id;

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