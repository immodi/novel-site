PRAGMA foreign_keys = ON;

-- Table: novels
CREATE TABLE IF NOT EXISTS novels (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    title TEXT NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    cover_image TEXT NOT NULL DEFAULT '',
    author TEXT NOT NULL DEFAULT '',
    status TEXT NOT NULL DEFAULT '',
    update_time TEXT NOT NULL DEFAULT '',
    latest_chapter_name TEXT NOT NULL DEFAULT '',
    total_chapters_number INTEGER NOT NULL DEFAULT 0
);

-- Table: chapters
CREATE TABLE IF NOT EXISTS chapters (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    novel_id INTEGER NOT NULL,
    title TEXT NOT NULL,
    content TEXT NOT NULL DEFAULT '',
    FOREIGN KEY (novel_id) REFERENCES novels(id) ON DELETE CASCADE
);

-- Index for faster lookups of chapters by novel
CREATE INDEX IF NOT EXISTS idx_chapters_novel_id ON chapters(novel_id);

-- Table: novel_genres (many-to-many since genres are a slice)
CREATE TABLE IF NOT EXISTS novel_genres (
    novel_id INTEGER NOT NULL,
    genre TEXT NOT NULL,
    PRIMARY KEY (novel_id, genre),
    FOREIGN KEY (novel_id) REFERENCES novels(id) ON DELETE CASCADE
);
