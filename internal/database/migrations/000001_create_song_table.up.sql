CREATE TABLE IF NOT EXISTS songs (
    id SERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    group_name TEXT NOT NULL,
    release_date TEXT NOT NULL,
    song_text TEXT NOT NULL,
    link TEXT NOT NULL
);