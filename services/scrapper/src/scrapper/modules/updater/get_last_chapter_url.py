import json
from pathlib import Path
from scrapper.helpers.helpers import CHAPTER_REGEX
from typing import Optional


def get_last_chapter_url(dir_path: str) -> Optional[str]:
    """
    Existing helper — unchanged.
    """
    dir_path_obj = Path(dir_path)
    json_entries: list[tuple[int, str]] = []

    for json_file in dir_path_obj.glob("*.json"):
        try:
            with open(json_file, "r", encoding="utf-8") as f:
                data = json.load(f)
            if isinstance(data, dict) and "url" in data:
                match = CHAPTER_REGEX.search(data["url"])
                if match:
                    chapter_num = int(match.group(1))
                    json_entries.append((chapter_num, data["url"]))
                else:
                    print(f"Skipping {json_file.name}: no chapter number in URL")
        except Exception as e:
            print(f"Skipping {json_file.name}: {e}")

    if not json_entries:
        return None

    json_entries.sort(key=lambda x: x[0])
    return json_entries[-1][1]
