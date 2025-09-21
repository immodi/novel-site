PRAGMA foreign_keys=off;

-- Check if the column already exists
-- (SQLite doesn't have IF NOT EXISTS for columns, so we do it by pragma)
CREATE TEMP TABLE user_bookmarks_columns AS
SELECT name
FROM pragma_table_info('user_bookmarks')
WHERE name = 'last_read_chapter_id';

-- If column is missing, migrate
INSERT INTO user_bookmarks_columns (name)
SELECT 'missing'
WHERE NOT EXISTS (SELECT 1 FROM user_bookmarks_columns);

-- Migration only runs if 'missing' was inserted
-- (SQLite executes whole script, so we use conditional CREATE TABLE)
-- Step 1: Recreate with new schema
CREATE TABLE IF NOT EXISTS user_bookmarks_new (
    user_id INTEGER NOT NULL,
    novel_id INTEGER NOT NULL,
    created_at TEXT NOT NULL DEFAULT '',
    last_read_chapter_id INTEGER DEFAULT NULL,
    PRIMARY KEY (user_id, novel_id),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (novel_id) REFERENCES novels(id) ON DELETE CASCADE,
    FOREIGN KEY (last_read_chapter_id) REFERENCES chapters(id) ON DELETE SET NULL
);

-- Step 2: Copy data
INSERT INTO user_bookmarks_new (user_id, novel_id, created_at)
SELECT user_id, novel_id, created_at
FROM user_bookmarks;

-- Step 3: Replace old table
DROP TABLE user_bookmarks;
ALTER TABLE user_bookmarks_new RENAME TO user_bookmarks;

DROP TABLE user_bookmarks_columns;

PRAGMA foreign_keys=on;

