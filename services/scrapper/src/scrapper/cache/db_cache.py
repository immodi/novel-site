import sqlite3
from os import makedirs, path
from scrapper.datatypes.novel import NovelData
from scrapper.datatypes.last_chapter_url_data import NovelLastChapter


class NovelDataCache:
    def __init__(self, db_path: str):
        makedirs(path.dirname(db_path), exist_ok=True)
        self.conn = sqlite3.connect(db_path)
        self.create_tables()

    def create_tables(self):
        with self.conn:
            self.conn.execute(
                """
                CREATE TABLE IF NOT EXISTS novels (
                    id INTEGER PRIMARY KEY AUTOINCREMENT,
                    title TEXT,
                    url TEXT UNIQUE
                )
                """
            )
            self.conn.execute(
                """
                CREATE TABLE IF NOT EXISTS last_novel_chapter (
                    id INTEGER PRIMARY KEY AUTOINCREMENT,
                    novel_id INTEGER UNIQUE,
                    chapter_url TEXT,
                    chapter_name TEXT,
                    FOREIGN KEY(novel_id) REFERENCES novels(id) ON DELETE CASCADE
                )
                """
            )

            self.conn.execute(
                """
                CREATE TABLE IF NOT EXISTS last_update (
                    id INTEGER PRIMARY KEY CHECK (id = 1),
                    minutes INTEGER
                )
                """
            )

    def get_all_novels(self) -> list[NovelData]:
        """Return all novels in the database as NovelData objects."""
        cursor = self.conn.cursor()
        cursor.execute("SELECT title, url FROM novels")
        rows = cursor.fetchall()
        novels = []
        for row in rows:
            novels.append(
                NovelData(
                    title=row[0],
                    author="",  # not stored
                    genres=[],  # not stored
                    status="",  # not stored
                    tags=[],  # not stored
                    cover_image="",  # not stored
                    description="",  # not stored
                    url=row[1],
                )
            )
        return novels

    def get_novel_by_url(self, url: str) -> tuple[int | None, NovelData | None]:
        """
        Get a novel by its URL.
        Returns a tuple: (novel_id, NovelData) or (None, None) if not found.
        """
        cursor = self.conn.cursor()
        cursor.execute("SELECT id, title, url FROM novels WHERE url = ?", (url,))
        row = cursor.fetchone()
        if row:
            novel_id = row[0]
            novel = NovelData(
                title=row[1],
                author="",  # not stored
                genres=[],  # not stored
                status="",  # not stored
                tags=[],  # not stored
                cover_image="",  # not stored
                description="",  # not stored
                url=row[2],
            )
            return novel_id, novel
        return None, None

    def save_novel(self, novel: NovelData):
        """Insert only if the URL does not already exist."""
        with self.conn:
            self.conn.execute(
                "INSERT OR IGNORE INTO novels (title, url) VALUES (?, ?)",
                (novel.title, novel.url),
            )

    def remove_novel_by_url(self, url: str):
        """Remove a novel by its URL (also removes its last chapter via cascade)."""
        with self.conn:
            self.conn.execute("DELETE FROM novels WHERE url = ?", (url,))

    def novel_exists(self, url: str) -> bool:
        cursor = self.conn.cursor()
        cursor.execute("SELECT 1 FROM novels WHERE url = ? LIMIT 1", (url,))
        return cursor.fetchone() is not None

    def get_novel_id(self, url: str) -> int | None:
        cursor = self.conn.cursor()
        cursor.execute("SELECT id FROM novels WHERE url = ?", (url,))
        row = cursor.fetchone()
        return row[0] if row else None

    def get_novel_url_by_chapter(self, chapter_url: str) -> str | None:
        """
        Get the novel URL corresponding to a given last chapter URL.
        Returns None if no matching novel is found.
        """
        cursor = self.conn.cursor()
        cursor.execute(
            """
            SELECT n.url
            FROM novels n
            JOIN last_novel_chapter lnc ON n.id = lnc.novel_id
            WHERE lnc.chapter_url = ?
            """,
            (chapter_url,),
        )
        row = cursor.fetchone()
        if row:
            return row[0]
        return None

    def save_last_chapter(self, novel_url: str, chapter_url: str, chapter_name: str):
        """Insert or update the last chapter for a novel (one-to-one)."""
        novel_id, novel = self.get_novel_by_url(novel_url)
        if novel is None or novel_id is None:
            raise ValueError("Novel must exist before adding a last chapter.")

        with self.conn:
            self.conn.execute(
                """
                INSERT INTO last_novel_chapter (novel_id, chapter_url, chapter_name)
                VALUES (?, ?, ?)
                ON CONFLICT(novel_id) DO UPDATE SET
                    chapter_url=excluded.chapter_url,
                    chapter_name=excluded.chapter_name
                """,
                (novel_id, chapter_url, chapter_name),
            )

    def get_last_chapter(self, novel_url: str) -> NovelLastChapter | None:
        """Get the last chapter of a novel as a NovelLastChapter object."""
        novel_id, novel = self.get_novel_by_url(novel_url)
        if novel is None or novel_id is None:
            return None

        cursor = self.conn.cursor()
        cursor.execute(
            "SELECT chapter_url, chapter_name FROM last_novel_chapter WHERE novel_id = ?",
            (novel_id,),
        )
        row = cursor.fetchone()
        if row:
            return NovelLastChapter(
                novel_name=novel.title,
                last_chapter_url=row[0],
                last_chapter_name=row[1],
            )
        return None

    def get_last_chapters(self, novels: list[NovelData]) -> list[NovelLastChapter]:
        """Get the last chapter of each novel in the list."""
        last_chapters = []
        for novel in novels:
            last_chapter = self.get_last_chapter(novel.url)
            if last_chapter is not None:
                last_chapters.append(last_chapter)
        return last_chapters

    def set_last_update(self, minutes: int):
        """Store a single integer (minutes) representing last update time."""
        with self.conn:
            self.conn.execute(
                """
                INSERT INTO last_update (id, minutes)
                VALUES (1, ?)
                ON CONFLICT(id) DO UPDATE SET minutes = excluded.minutes
                """,
                (minutes,),
            )

    def get_last_update(self) -> int | None:
        """Return the stored minutes value, or None if not set."""
        cursor = self.conn.cursor()
        cursor.execute("SELECT minutes FROM last_update WHERE id = 1")
        row = cursor.fetchone()
        return row[0] if row else None

    def close(self):
        self.conn.close()
