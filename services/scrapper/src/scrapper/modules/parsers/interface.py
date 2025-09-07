from typing import Protocol, List
from lxml import html
from scrapper import config
from scrapper.helpers import utils
from scrapper.datatypes.novel import (
    ChapterData,
    NovelData,
    NovelLink,
)
from os import path
from enum import Enum


class SkipDuplicate(Enum):
    NONE = "none"
    NOVEL = "novel"
    CHAPTER = "chapter"


class Parser(Protocol):
    def parse_list_of_novels(self, tree: html.HtmlElement) -> List[NovelLink]: ...
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
