CREATE TABLE IF NOT EXISTS groups(
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS songs (
    id SERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    group_id INTEGER NOT NULL,
    release_date DATE NOT NULL,
    song_text TEXT NOT NULL,
    link TEXT NOT NULL
);