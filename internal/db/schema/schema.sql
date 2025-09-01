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
    view_count INTEGER NOT NULL DEFAULT 0
);

-- Table: chapters (enhanced with chapter_number)
CREATE TABLE IF NOT EXISTS chapters (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    novel_id INTEGER NOT NULL,
    chapter_number INTEGER NOT NULL,
    title TEXT NOT NULL,
    content TEXT NOT NULL DEFAULT '',
    FOREIGN KEY (novel_id) REFERENCES novels(id) ON DELETE CASCADE,
    UNIQUE(novel_id, chapter_number) -- Ensure unique chapter numbers per novel
);

-- Index for faster lookups of chapters by novel
CREATE INDEX IF NOT EXISTS idx_chapters_novel_id ON chapters(novel_id);

-- Index for faster lookups of chapters by novel and chapter number
CREATE INDEX IF NOT EXISTS idx_chapters_novel_chapter ON chapters(novel_id, chapter_number);

-- Table: novel_genres (many-to-many since genres are a slice)
CREATE TABLE IF NOT EXISTS novel_genres (
    novel_id INTEGER NOT NULL,
    genre TEXT NOT NULL,
    PRIMARY KEY (novel_id, genre),
    FOREIGN KEY (novel_id) REFERENCES novels(id) ON DELETE CASCADE
);


-- Table: novel_tags (many-to-many for tags)
CREATE TABLE IF NOT EXISTS novel_tags (
    novel_id INTEGER NOT NULL,
    tag TEXT NOT NULL,
    PRIMARY KEY (novel_id, tag),
    FOREIGN KEY (novel_id) REFERENCES novels(id) ON DELETE CASCADE
);

-- Optional: index for faster lookups by novel
CREATE INDEX IF NOT EXISTS idx_novel_tags_novel_id ON novel_tags(novel_id);

-- Optional: index for faster lookups by tag (if you want to query novels by tag)
CREATE INDEX IF NOT EXISTS idx_novel_tags_tag ON novel_tags(tag);

