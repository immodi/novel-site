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
