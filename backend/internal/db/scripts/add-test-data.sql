DO $$
DECLARE inserted_venue INTEGER;
BEGIN
    -- insert venue, store ID in variable 'inserted_venue'
    INSERT INTO london_jam_sessions.venues (
        venue_name, address_first_line, city, postcode, geom, venue_website, backline, venue_comments
    ) VALUES (
        'The Spice of Life', '6 Moor Street', 'London', 'W1D 5NA', 'SRID=4326;POINT(-0.12980699914457172 51.51339645)', 'https://www.spiceoflifesoho.com', '{PA,Drums,Bass_Amp,Guitar_Amp}', '{Sessions in basement}'
    ) RETURNING venue_id INTO inserted_venue;

    -- insert sessions, use 'inserted_venue' as fkey
    INSERT INTO london_jam_sessions.jamsessions (
        session_name, venue, description, start_time_utc, duration_minutes, interval, session_website, genres, session_comments
    ) VALUES (
        'Jazz comments Jazz Jam', inserted_venue, 'House band plays first set. Free entry.', '2024-08-25T14:00:00Z', 180, 'Weekly', 'https://www.spiceoflifesoho.com/events/different-planet-presents-jazz-comments-jazz-jam-6/', '{Straight-Ahead_Jazz}', '{House band plays 1st hour,Sign up for jam at door}'
    );

    INSERT INTO london_jam_sessions.jamsessions (
        session_name, venue, description, start_time_utc, duration_minutes, interval, session_website, genres, session_comments
    ) VALUES (
        'Daily Jam Placeholder', inserted_venue, 'House band plays first set. Free entry.', '2024-08-25T19:00:00Z', 120, 'Daily', 'https://www.spiceoflifesoho.com/events/different-planet-presents-jazz-comments-jazz-jam-6/', '{Straight-Ahead_Jazz,Blues}', '{Lorem ipsum dolor sit amet,Consectetur adipiscing elit,Sed do eiusmod tempor incididunt ut labore et dolore magna aliqua}'
    );
END $$;