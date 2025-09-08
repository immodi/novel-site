import os
import json
from typing import List, Union
from scrapper import config
from scrapper.datatypes.novel import ChapterData, NovelData
from dataclasses import asdict, is_dataclass
from scrapper.helpers import utils


def to_serializable(obj):
    """Convert dataclass instances (not classes) to dicts recursively."""
    if is_dataclass(obj) and not isinstance(obj, type):
        return asdict(obj)
    if isinstance(obj, list):
        return [to_serializable(x) for x in obj]
    if isinstance(obj, dict):
        return {k: to_serializable(v) for k, v in obj.items()}
    return obj


def save_item(novel_data: Union[NovelData, ChapterData], directory: str) -> str:
    os.makedirs(directory, exist_ok=True)

    title = novel_data.title or "untitled"
    slug = utils.slugify(title)
    filepath = os.path.join(directory, f"{slug}.json")

    with open(filepath, "w", encoding="utf-8") as f:
        json.dump(to_serializable(novel_data), f, ensure_ascii=False, indent=2)

    return filepath


def save_all(
    novels: Union[List[NovelData], List[ChapterData]], filename: str = "novels.json"
) -> str:
    os.makedirs(config.OUTPUT_DIR, exist_ok=True)
    filepath = os.path.join(config.OUTPUT_DIR, filename)

    with open(filepath, "w", encoding="utf-8") as f:
        json.dump(to_serializable(novels), f, ensure_ascii=False, indent=2)

    return filepath
