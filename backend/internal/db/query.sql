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
    SELECT s.*, l.*, coalesce(round(avg(rating), 2), 0.0)::real AS rating FROM london_jam_sessions.jamsessions s
    JOIN london_jam_sessions.venues l ON s.venue = l.venue_id
    LEFT OUTER JOIN london_jam_sessions.ratings r ON s.session_id = r.session
    WHERE interval = 'Daily' 
    OR (
        start_time_utc::date = sqlc.arg(date)::date
        AND interval = 'Once'
    ) 
    OR (
        interval = 'Weekly' AND EXTRACT(dow FROM s.start_time_utc) = EXTRACT(dow FROM sqlc.arg(date)::date)
    ) 
    OR (
        interval = 'Fortnightly' AND (sqlc.arg(date)::date - start_time_utc::date) % 14 = 0 -- check if the number of days between is divisible by 14
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
        EXTRACT(MONTH FROM (sqlc.arg(date)::date + interval '7 day')) != EXTRACT(MONTH FROM sqlc.arg(date)) 
    )
    GROUP BY s.session_id, l.venue_id
)
SELECT json_build_object(
    'type', 'FeatureCollection',
    'features', json_agg(public.ST_AsGeoJSON(t.*)::json)
) FROM t;

-- name: GetSessionsByDateRangeAsGeoJSON :one
WITH dates_in_range AS (
    SELECT generate_series(sqlc.arg(start_date)::date, sqlc.arg(end_date)::date, interval '1 day') d
), t AS (
    SELECT s.*, l.*, coalesce(round(avg(rating), 2), 0.0)::real AS rating FROM london_jam_sessions.jamsessions s
    JOIN london_jam_sessions.venues l ON s.venue = l.venue_id
    LEFT OUTER JOIN london_jam_sessions.ratings r ON s.session_id = r.session
    WHERE interval = 'Daily' 
    OR (
        start_time_utc::date >= sqlc.arg(start_date)::date
        AND start_time_utc::date < sqlc.arg(end_date)::date + interval '1 day'
        AND interval = 'Once'
    ) 
    OR (
        interval = 'Weekly' 
        AND EXTRACT(dow FROM sqlc.arg(start_date)::date) <= EXTRACT(dow FROM s.start_time_utc)
        AND EXTRACT(dow FROM sqlc.arg(end_date)::date) >= EXTRACT(dow FROM s.start_time_utc)
    ) 
    OR (
        interval = 'Fortnightly' AND EXISTS (
            SELECT *
            FROM dates_in_range
            WHERE (dates_in_range.d::date - s.start_time_utc::date) % 14 = 0 
            -- generate a list of dates between start/end
            -- check if for any of them the condition is true
        )
    )
    OR (
        interval = 'FirstOfMonth' AND EXISTS (
            SELECT *
            FROM dates_in_range
            WHERE CEIL (EXTRACT(DAY FROM dates_in_range.d) / 7) = 1
        ) 
    )
    OR (
        interval = 'SecondOfMonth' AND EXISTS (
            SELECT *
            FROM dates_in_range
            WHERE CEIL (EXTRACT(DAY FROM dates_in_range.d) / 7) = 2
        ) 
    )
    OR (
        interval = 'ThirdOfMonth' AND EXISTS (
            SELECT *
            FROM dates_in_range
            WHERE CEIL (EXTRACT(DAY FROM dates_in_range.d) / 7) = 3
        ) 
    )
    OR (
        interval = 'FourthOfMonth' AND EXISTS (
            SELECT *
            FROM dates_in_range
            WHERE CEIL (EXTRACT(DAY FROM dates_in_range.d) / 7) = 4
        ) 
    )
    OR interval = 'LastOfMonth' AND EXISTS (
        SELECT *
        FROM dates_in_range
        WHERE EXTRACT(MONTH FROM (dates_in_range.d + interval '7 day')) != EXTRACT(MONTH FROM dates_in_range.d) 
    )
    GROUP BY s.session_id, l.venue_id
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
    JOIN london_jam_sessions.venues l ON s.venue = l.venue_id
    LEFT OUTER JOIN london_jam_sessions.ratings r ON s.session_id = r.session
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
        interval = 'Fortnightly' AND (sqlc.arg(date)::date - start_time_utc::date) % 14 = 0
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
        EXTRACT(MONTH FROM (sqlc.arg(date)::date + interval '7 day')) != EXTRACT(MONTH FROM sqlc.arg(date)::date) 
    )
    GROUP BY s.session_id, l.venue_id
)
SELECT json_build_object(
    'type', 'FeatureCollection',
    'features', json_agg(public.ST_AsGeoJSON(t.*)::json)
) FROM t;

-- name: GetSessionsByDateRangeAndGenreAsGeoJSON :one
WITH dates_in_range AS (
    SELECT generate_series(sqlc.arg(start_date)::date, sqlc.arg(end_date)::date, interval '1 day') d
), t AS (
    SELECT s.*, l.*, coalesce(round(avg(rating), 2), 0.0)::real AS rating FROM london_jam_sessions.jamsessions s
    JOIN london_jam_sessions.venues l ON s.venue = l.venue_id
    LEFT OUTER JOIN london_jam_sessions.ratings r ON s.session_id = r.session
    WHERE s.genres @> sqlc.arg(genres) 
    AND interval = 'Daily'
    OR (
        start_time_utc::date >= sqlc.arg(start_date)::date
        AND start_time_utc::date < sqlc.arg(end_date)::date + interval '1 day'
        AND interval = 'Once'
    ) 
    OR (
        interval = 'Weekly' 
        AND EXTRACT(dow FROM sqlc.arg(start_date)::date) <= EXTRACT(dow FROM s.start_time_utc)
        AND EXTRACT(dow FROM sqlc.arg(end_date)::date) >= EXTRACT(dow FROM s.start_time_utc)
    ) 
    OR (
        interval = 'Fortnightly' AND EXISTS (
            SELECT *
            FROM dates_in_range
            WHERE (dates_in_range.d::date - s.start_time_utc::date) % 14 = 0 
            -- generate a list of dates between start/end
            -- check if for any of them the condition is true
        )
    )
    OR (
        interval = 'FirstOfMonth' AND EXISTS (
            SELECT *
            FROM dates_in_range
            WHERE CEIL (EXTRACT(DAY FROM dates_in_range.d) / 7) = 1
        ) 
    )
    OR (
        interval = 'SecondOfMonth' AND EXISTS (
            SELECT *
            FROM dates_in_range
            WHERE CEIL (EXTRACT(DAY FROM dates_in_range.d) / 7) = 2
        ) 
    )
    OR (
        interval = 'ThirdOfMonth' AND EXISTS (
            SELECT *
            FROM dates_in_range
            WHERE CEIL (EXTRACT(DAY FROM dates_in_range.d) / 7) = 3
        ) 
    )
    OR (
        interval = 'FourthOfMonth' AND EXISTS (
            SELECT *
            FROM dates_in_range
            WHERE CEIL (EXTRACT(DAY FROM dates_in_range.d) / 7) = 4
        ) 
    )
    OR interval = 'LastOfMonth' AND EXISTS (
        SELECT *
        FROM dates_in_range
        WHERE EXTRACT(MONTH FROM (dates_in_range.d + interval '7 day')) != EXTRACT(MONTH FROM dates_in_range.d) 
    )
    GROUP BY s.session_id, l.venue_id
)
SELECT json_build_object(
    'type', 'FeatureCollection',
    'features', json_agg(public.ST_AsGeoJSON(t.*)::json)
) FROM t;


-- name: GetSessionsByDateAndBacklineAsGeoJSON :one
WITH t AS (
    SELECT s.*, l.*, coalesce(round(avg(rating), 2), 0.0)::real AS rating FROM london_jam_sessions.jamsessions s
    JOIN london_jam_sessions.venues l ON s.venue = l.venue_id
    LEFT OUTER JOIN london_jam_sessions.ratings r ON s.session_id = r.session
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
        interval = 'Fortnightly' AND (sqlc.arg(date)::date - start_time_utc::date) % 14 = 0
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
        EXTRACT(MONTH FROM (sqlc.arg(date)::date + interval '7 day')) != EXTRACT(MONTH FROM sqlc.arg(date)::date) 
    )
    GROUP BY s.session_id, l.venue_id
)
SELECT json_build_object(
    'type', 'FeatureCollection',
    'features', json_agg(public.ST_AsGeoJSON(t.*)::json)
) FROM t;

-- name: GetSessionsByDateRangeAndBacklineAsGeoJSON :one
WITH dates_in_range AS (
    SELECT generate_series(sqlc.arg(start_date)::date, sqlc.arg(end_date)::date, interval '1 day') d
), t AS (
    SELECT s.*, l.*, coalesce(round(avg(rating), 2), 0.0)::real AS rating FROM london_jam_sessions.jamsessions s
    JOIN london_jam_sessions.venues l ON s.venue = l.venue_id
    LEFT OUTER JOIN london_jam_sessions.ratings r ON s.session_id = r.session
    WHERE l.backline @> sqlc.arg(backline)
    AND interval = 'Daily'
    OR (
        start_time_utc::date >= sqlc.arg(start_date)::date
        AND start_time_utc::date < sqlc.arg(end_date)::date + interval '1 day'
        AND interval = 'Once'
    ) 
    OR (
        interval = 'Weekly' 
        AND EXTRACT(dow FROM sqlc.arg(start_date)::date) <= EXTRACT(dow FROM s.start_time_utc)
        AND EXTRACT(dow FROM sqlc.arg(end_date)::date) >= EXTRACT(dow FROM s.start_time_utc)
    ) 
    OR (
        interval = 'Fortnightly' AND EXISTS (
            SELECT *
            FROM dates_in_range
            WHERE (dates_in_range.d::date - s.start_time_utc::date) % 14 = 0 
            -- generate a list of dates between start/end
            -- check if for any of them the condition is true
        )
    )
    OR (
        interval = 'FirstOfMonth' AND EXISTS (
            SELECT *
            FROM dates_in_range
            WHERE CEIL (EXTRACT(DAY FROM dates_in_range.d) / 7) = 1
        ) 
    )
    OR (
        interval = 'SecondOfMonth' AND EXISTS (
            SELECT *
            FROM dates_in_range
            WHERE CEIL (EXTRACT(DAY FROM dates_in_range.d) / 7) = 2
        ) 
    )
    OR (
        interval = 'ThirdOfMonth' AND EXISTS (
            SELECT *
            FROM dates_in_range
            WHERE CEIL (EXTRACT(DAY FROM dates_in_range.d) / 7) = 3
        ) 
    )
    OR (
        interval = 'FourthOfMonth' AND EXISTS (
            SELECT *
            FROM dates_in_range
            WHERE CEIL (EXTRACT(DAY FROM dates_in_range.d) / 7) = 4
        ) 
    )
    OR interval = 'LastOfMonth' AND EXISTS (
        SELECT *
        FROM dates_in_range
        WHERE EXTRACT(MONTH FROM (dates_in_range.d + interval '7 day')) != EXTRACT(MONTH FROM dates_in_range.d) 
    )
    GROUP BY s.session_id, l.venue_id
)
SELECT json_build_object(
    'type', 'FeatureCollection',
    'features', json_agg(public.ST_AsGeoJSON(t.*)::json)
) FROM t;

-- name: GetSessionsByDateAndGenreAndBacklineAsGeoJSON :one
WITH t AS (
    SELECT s.*, l.*, coalesce(round(avg(rating), 2), 0.0)::real AS rating FROM london_jam_sessions.jamsessions s
    JOIN london_jam_sessions.venues l ON s.venue = l.venue_id
    LEFT OUTER JOIN london_jam_sessions.ratings r ON s.session_id = r.session
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
        interval = 'Fortnightly' AND (sqlc.arg(date)::date - start_time_utc::date) % 14 = 0
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
        EXTRACT(MONTH FROM (sqlc.arg(date)::date + interval '7 day')) != EXTRACT(MONTH FROM sqlc.arg(date)::date) 
    )
    GROUP BY s.session_id,l.venue_id
)
SELECT json_build_object(
    'type', 'FeatureCollection',
    'features', json_agg(public.ST_AsGeoJSON(t.*)::json)
) FROM t;

-- name: GetSessionsByDateRangeAndGenreAndBacklineAsGeoJSON :one
WITH dates_in_range AS (
    SELECT generate_series(sqlc.arg(start_date)::date, sqlc.arg(end_date)::date, interval '1 day') d
), t AS (
    SELECT s.*, l.*, coalesce(round(avg(rating), 2), 0.0)::real AS rating FROM london_jam_sessions.jamsessions s
    JOIN london_jam_sessions.venues l ON s.venue = l.venue_id
    LEFT OUTER JOIN london_jam_sessions.ratings r ON s.session_id = r.session
    WHERE s.genres @> sqlc.arg(genres)
    AND l.backline @> sqlc.arg(backline)
    AND interval = 'Daily'
    OR (
        start_time_utc::date >= sqlc.arg(start_date)::date
        AND start_time_utc::date < sqlc.arg(end_date)::date + interval '1 day'
        AND interval = 'Once'
    ) 
    OR (
        interval = 'Weekly' 
        AND EXTRACT(dow FROM sqlc.arg(start_date)::date) <= EXTRACT(dow FROM s.start_time_utc)
        AND EXTRACT(dow FROM sqlc.arg(end_date)::date) >= EXTRACT(dow FROM s.start_time_utc)
    ) 
    OR (
        interval = 'Fortnightly' AND EXISTS (
            SELECT *
            FROM dates_in_range
            WHERE (dates_in_range.d::date - s.start_time_utc::date) % 14 = 0 
            -- generate a list of dates between start/end
            -- check if for any of them the condition is true
        )
    )
    OR (
        interval = 'FirstOfMonth' AND EXISTS (
            SELECT *
            FROM dates_in_range
            WHERE CEIL (EXTRACT(DAY FROM dates_in_range.d) / 7) = 1
        ) 
    )
    OR (
        interval = 'SecondOfMonth' AND EXISTS (
            SELECT *
            FROM dates_in_range
            WHERE CEIL (EXTRACT(DAY FROM dates_in_range.d) / 7) = 2
        ) 
    )
    OR (
        interval = 'ThirdOfMonth' AND EXISTS (
            SELECT *
            FROM dates_in_range
            WHERE CEIL (EXTRACT(DAY FROM dates_in_range.d) / 7) = 3
        ) 
    )
    OR (
        interval = 'FourthOfMonth' AND EXISTS (
            SELECT *
            FROM dates_in_range
            WHERE CEIL (EXTRACT(DAY FROM dates_in_range.d) / 7) = 4
        ) 
    )
    OR interval = 'LastOfMonth' AND EXISTS (
        SELECT *
        FROM dates_in_range
        WHERE EXTRACT(MONTH FROM (dates_in_range.d + interval '7 day')) != EXTRACT(MONTH FROM dates_in_range.d) 
    )
    GROUP BY s.session_id, l.venue_id
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