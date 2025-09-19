BASE_URL = "https://novelfire.net"
LIST_URL = f"{BASE_URL}/genre-all/sort-popular/status-all/all-novel"
OUTPUT_DIR = "data"
CACHE_DB_PATH = f"{OUTPUT_DIR}/cache.db"
HEADERS = {
    "User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) "
    "AppleWebKit/537.36 (KHTML, like Gecko) "
    "Chrome/120.0.0.0 Safari/537.36",
    "Referer": BASE_URL,
}
