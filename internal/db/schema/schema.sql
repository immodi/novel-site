PRAGMA foreign_keys = ON;

-- Table: novels
CREATE TABLE IF NOT EXISTS novels (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    title TEXT NOT NULL,
    slug TEXT NOT NULL, -- URL-friendly, unique identifier
    description TEXT NOT NULL DEFAULT '',
    cover_image TEXT NOT NULL DEFAULT '',
    author TEXT NOT NULL,
    author_slug TEXT NOT NULL, -- URL-friendly author
    publisher TEXT NOT NULL DEFAULT '',
    release_year INTEGER NOT NULL DEFAULT 0,
    is_completed INTEGER NOT NULL DEFAULT 0, -- 0 = ongoing, 1 = completed
    update_time TEXT NOT NULL DEFAULT '',
    view_count INTEGER NOT NULL DEFAULT 0
);

-- Unique indexes for fast and safe lookups
CREATE UNIQUE INDEX IF NOT EXISTS idx_novels_title_lower ON novels (LOWER(title));
CREATE UNIQUE INDEX IF NOT EXISTS idx_novels_slug_lower ON novels (LOWER(slug));
CREATE INDEX IF NOT EXISTS idx_novels_author_slug_lower ON novels (LOWER(author_slug));

-- Table: chapters
CREATE TABLE IF NOT EXISTS chapters (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    novel_id INTEGER NOT NULL,
    chapter_number INTEGER NOT NULL,
    title TEXT NOT NULL,
    content TEXT NOT NULL DEFAULT '',
    FOREIGN KEY (novel_id) REFERENCES novels(id) ON DELETE CASCADE,
    UNIQUE(novel_id, chapter_number)
);

-- Indexes for faster lookups of chapters
CREATE INDEX IF NOT EXISTS idx_chapters_novel_id ON chapters(novel_id);
CREATE INDEX IF NOT EXISTS idx_chapters_novel_chapter ON chapters(novel_id, chapter_number);

-- Table: novel_genres (many-to-many, free-form)
CREATE TABLE IF NOT EXISTS novel_genres (
    novel_id INTEGER NOT NULL,
    genre TEXT NOT NULL,
    genre_slug TEXT NOT NULL, -- URL-friendly version
    PRIMARY KEY (novel_id, genre),
    FOREIGN KEY (novel_id) REFERENCES novels(id) ON DELETE CASCADE
);

-- Table: novel_tags (many-to-many, free-form)
CREATE TABLE IF NOT EXISTS novel_tags (
    novel_id INTEGER NOT NULL,
    tag TEXT NOT NULL,
    tag_slug TEXT NOT NULL, -- URL-friendly version
    PRIMARY KEY (novel_id, tag),
    FOREIGN KEY (novel_id) REFERENCES novels(id) ON DELETE CASCADE
);

-- Indexes for faster lookups
CREATE INDEX IF NOT EXISTS idx_novel_tags_novel_id ON novel_tags(novel_id);
CREATE INDEX IF NOT EXISTS idx_novel_tags_tag_slug ON novel_tags(tag_slug);
CREATE INDEX IF NOT EXISTS idx_novel_genres_novel_id ON novel_genres(novel_id);
CREATE INDEX IF NOT EXISTS idx_novel_genres_genre_slug ON novel_genres(genre_slug);

-- Table: users
CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    username TEXT NOT NULL UNIQUE,
    email TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL, 
    role TEXT NOT NULL DEFAULT '', 
    created_at TEXT NOT NULL DEFAULT '',
    image TEXT NOT NULL DEFAULT 'https://www.citypng.com/public/uploads/preview/white-user-member-guest-icon-png-image-701751695037005zdurfaim0y.png' -- profile image stored as base64 
);

-- Indexes for quick lookups
CREATE UNIQUE INDEX IF NOT EXISTS idx_users_username_lower ON users (LOWER(username));
CREATE UNIQUE INDEX IF NOT EXISTS idx_users_email_lower ON users (LOWER(email));

-- Table: user_bookmarks (many-to-many between users and novels)
CREATE TABLE IF NOT EXISTS user_bookmarks (
    user_id INTEGER NOT NULL,
    novel_id INTEGER NOT NULL,
    created_at TEXT NOT NULL DEFAULT '',
    PRIMARY KEY (user_id, novel_id),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (novel_id) REFERENCES novels(id) ON DELETE CASCADE
);

-- Indexes for bookmarks
CREATE INDEX IF NOT EXISTS idx_user_bookmarks_user_id ON user_bookmarks(user_id);
CREATE INDEX IF NOT EXISTS idx_user_bookmarks_novel_id ON user_bookmarks(novel_id);
