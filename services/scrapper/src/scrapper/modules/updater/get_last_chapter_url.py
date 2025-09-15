import json
from pathlib import Path
from scrapper.helpers.helpers import CHAPTER_REGEX
from typing import Optional, Tuple


def get_last_chapter_data(dir_path: str) -> Optional[Tuple[str, str]]:
    """
    Get the URL and title of the latest chapter JSON in the given directory.

    Parameters
    ----------
    dir_path : str
        Path to the directory containing chapter JSON files.

    Returns
    -------
    Optional[Tuple[str, str]]
        (url, title) of the last chapter found, or None if none found.
    """
    dir_path_obj = Path(dir_path)
    json_entries: list[tuple[int, str, str]] = []  # (chapter_num, url, title)

    for json_file in dir_path_obj.glob("*.json"):
        try:
            with open(json_file, "r", encoding="utf-8") as f:
                data = json.load(f)

            if isinstance(data, dict) and "url" in data and "title" in data:
                match = CHAPTER_REGEX.search(data["url"])
                if match:
                    chapter_num = int(match.group(1))
                    json_entries.append((chapter_num, data["url"], data["title"]))
                else:
                    print(f"Skipping {json_file.name}: no chapter number in URL")
            else:
                print(f"Skipping {json_file.name}: missing 'url' or 'title'")
        except Exception as e:
            print(f"Skipping {json_file.name}: {e}")

    if not json_entries:
        return None

    # sort by chapter number and return (url, title) of the latest
    json_entries.sort(key=lambda x: x[0])
    _, url, title = json_entries[-1]
    return url, title
