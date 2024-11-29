DO $$
DECLARE inserted_venue INTEGER;
DECLARE inserted_session INTEGER;
DECLARE inserted_comment INTEGER;
BEGIN
    -- insert venue, store ID in variable 'inserted_venue'
    INSERT INTO london_jam_sessions.venues (
        venue_name, address_first_line, city, postcode, geom, venue_website, backline, venue_comments
    ) VALUES (
        'The Spice of Life', '6 Moor Street', 'London', 'W1D 5NA', 'SRID=4326;POINT(-0.12980699914457172 51.51339645)', 'https://www.spiceoflifesoho.com', '{PA,Drums,Bass_Amp,Guitar_Amp}', '{Sessions in basement}'
    ) RETURNING venue_id INTO inserted_venue;

    -- insert sessions, use 'inserted_venue' as fkey
    INSERT INTO london_jam_sessions.jamsessions (
        session_name, venue, description, start_time_utc, duration_minutes, interval, session_website, genres
    ) VALUES (
        'Jazz Jam', inserted_venue, 'House band plays first set. Free entry.', '2024-08-25T14:00:00Z', 180, 'Weekly', 'https://www.spiceoflifesoho.com/events/different-planet-presents-jazz-comments-jazz-jam-6/', '{Straight-Ahead_Jazz}'
    );

    INSERT INTO london_jam_sessions.jamsessions (
        session_name, venue, description, start_time_utc, duration_minutes, interval, session_website, genres
    ) VALUES (
        'Irregular Jam', inserted_venue, 'House band plays first set. Free entry.', '2024-08-25T14:00:00Z', 180, 'IrregularWeekly', 'https://www.spiceoflifesoho.com/events/different-planet-presents-jazz-comments-jazz-jam-6/', '{Straight-Ahead_Jazz}'
    );

    INSERT INTO london_jam_sessions.jamsessions (
        session_name, venue, description, start_time_utc, duration_minutes, interval, session_website, genres
    ) VALUES (
        'Daily Jam Placeholder', inserted_venue, 'House band plays first set. Free entry.', '2024-08-25T19:00:00Z', 120, 'Daily', 'https://www.spiceoflifesoho.com/events/different-planet-presents-jazz-comments-jazz-jam-6/', '{Straight-Ahead_Jazz,Blues}'
    ) RETURNING session_id INTO inserted_session;


    -- add comments
    INSERT INTO london_jam_sessions.comments (
        author, content, session
    ) VALUES (
        'Jazz cat', 'House band plays 1st hour, sign up for jam at door.', inserted_session
    ) RETURNING comment_id INTO inserted_comment;

    INSERT INTO london_jam_sessions.comments (
        author, content, session
    ) VALUES (
        'John Doe', 'Lorem ipsum dolor sit amet, consectetur adipiscing elit. Etiam euismod diam eget felis imperdiet sollicitudin. Quisque sagittis justo diam, vitae rutrum ligula ornare at. Ut volutpat eleifend quam et malesuada.', inserted_session
    ) RETURNING comment_id INTO inserted_comment;

    INSERT INTO london_jam_sessions.ratings (
        session, rating, comment
    ) VALUES (
        inserted_session, 3, inserted_comment
    );
END $$;