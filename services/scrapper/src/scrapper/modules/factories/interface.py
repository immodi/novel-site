from typing import Protocol, List, Union
from lxml import html
from scrapper import config
from scrapper.helpers import utils
from scrapper.datatypes.novel import (
    ChapterData,
    NovelData,
    NovelLink,
)
from os import path, remove
from enum import Enum


class SkipDuplicate(Enum):
    NONE = "none"
    NOVEL = "novel"
    CHAPTER = "chapter"


class Parser(Protocol):
    def parse_list_of_novels(
        self, tree: Union[html.HtmlElement, List[html.HtmlElement]]
    ) -> List[NovelLink]: ...
    def parse_novel(
        self, tree: html.HtmlElement, url: str, save_image: bool = True
    ) -> NovelData: ...
    def parse_chapters(
        self, url: str, novel_name: str, save_per_chapter: bool
    ) -> List[ChapterData]: ...

    def novel_exists(self, title: str) -> bool:
        """Check if a novel JSON file already exists in DATA_DIR"""
        safe_title = utils.slugify(title)
        file_path = path.join(config.OUTPUT_DIR, "novels", f"{safe_title}.json")
        return path.exists(file_path)

    def chapter_exists(self, title: str, novel_name: str) -> bool:
        """Check if a chapter JSON file already exists in DATA_DIR"""
        safe_title = utils.slugify(title)
        file_path = path.join(
            config.OUTPUT_DIR, "chapters", novel_name, f"{safe_title}.json"
        )

        return path.exists(file_path)

    def clean_up_novel(self, novel_name: str):
        json_file = f"{config.OUTPUT_DIR}/novels/{utils.slugify(novel_name)}.json"
        cover_file = f"{config.OUTPUT_DIR}/covers/{utils.slugify(novel_name)}.jpg"

        for file_path in [json_file, cover_file]:
            try:
                remove(file_path)
                print(f"Removed {file_path}")
            except FileNotFoundError:
                print(f"File not found: {file_path}")
            except Exception as e:
                print(f"Failed to remove {file_path}: {e}")
