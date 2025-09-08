import os
import requests
from urllib.parse import urlparse
import json
from pathlib import Path
from typing import Optional, Callable, Tuple
import tempfile
from scrapper.helpers import utils
import re

CHAPTER_REGEX = re.compile(r"chapter-(\d+)", re.IGNORECASE)


def download_image(
    src: str, directory: str, image_name: Optional[str] = None
) -> Optional[str]:
    try:
        os.makedirs(directory, exist_ok=True)

        # Extract extension from URL
        url_path = urlparse(src).path
        ext = os.path.splitext(url_path)[1] or ".jpg"

        # Decide filename
        if image_name:
            filename = f"{utils.slugify(image_name)}{ext}"
        else:
            filename = os.path.basename(url_path) or f"image{ext}"

        filepath = os.path.join(directory, filename)

        headers = {
            "User-Agent": (
                "Mozilla/5.0 (X11; Linux x86_64) "
                "AppleWebKit/537.36 (KHTML, like Gecko) "
                "Chrome/127.0.0.0 Safari/537.36"
            ),
            "Referer": "https://novelfire.net/",
        }

        response = requests.get(src, headers=headers, timeout=10)
        response.raise_for_status()

        with open(filepath, "wb") as f:
            f.write(response.content)

        return filepath
    except Exception as e:
        print(f"Failed to download image: {e}")
        return None


def combine_json_objects_to_array(dir_path: str) -> Tuple[str, Callable[[], None]]:
    dir_path_obj = Path(dir_path)  # convert to Path
    combined = []

    # Collect JSON objects and pair with their chapter number
    json_entries = []
    for json_file in dir_path_obj.glob("*.json"):
        try:
            with open(json_file, "r", encoding="utf-8") as f:
                data = json.load(f)
                if isinstance(data, dict) and "url" in data:
                    match = CHAPTER_REGEX.search(data["url"])
                    if match:
                        chapter_num = int(match.group(1))
                        json_entries.append((chapter_num, data))
                    else:
                        print(
                            f"Skipping {json_file.name}: could not find chapter number in URL"
                        )
        except Exception as e:
            print(f"Skipping {json_file.name}: {e}")

    # Sort by chapter number
    json_entries.sort(key=lambda x: x[0])

    # Extract only the data part
    combined = [entry[1] for entry in json_entries]

    # Create a temp file in the system temp dir with the last dir name
    temp_dir = Path(tempfile.gettempdir())
    temp_file_path = temp_dir / f"{dir_path_obj.name}.json"

    with open(temp_file_path, "w", encoding="utf-8") as f:
        json.dump(combined, f)

    print(f"Combined {len(combined)} JSON objects into {temp_file_path}")

    # Cleanup callback
    def cleanup():
        if temp_file_path.exists():
            temp_file_path.unlink()
            print(f"Deleted temp file: {temp_file_path}")

    return str(temp_file_path), cleanup
