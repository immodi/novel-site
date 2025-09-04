import os
import json
from typing import Any, Mapping, Sequence
from modules import config, utils


def save_item(novel_data: Mapping[str, Any], directory: str = config.OUTPUT_DIR) -> str:
    os.makedirs(directory, exist_ok=True)  # <--- use the passed directory

    title = novel_data.get("title") or "untitled"
    slug = utils.slugify(title)
    filepath = os.path.join(directory, f"{slug}.json")
    with open(filepath, "w+", encoding="utf-8") as f:
        json.dump(novel_data, f, ensure_ascii=False, indent=2)
    return filepath


def save_all(novels: Sequence[Mapping[str, Any]], filename: str = "novels.json") -> str:
    os.makedirs(config.OUTPUT_DIR, exist_ok=True)
    filepath = os.path.join(config.OUTPUT_DIR, filename)
    with open(filepath, "w", encoding="utf-8") as f:
        json.dump(novels, f, ensure_ascii=False, indent=2)
    return filepath
