from typing import Protocol, List, Union, Tuple, Generator
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
    """
    takes either a list OR one html document(s) and parsers them into a list of links for each individual novel page
    """

    def parse_list_of_novels(
        self, tree: Union[html.HtmlElement, List[html.HtmlElement]]
    ) -> List[NovelLink]: ...

    """
    takes an html document of the web novel and extract the data
    """

    def parse_novel(
        self, tree: html.HtmlElement, url: str, save_image: bool = True
    ) -> NovelData: ...

    """
    Fetch chapters using the first link and navigating via #next_chap button
    """

    def parse_chapters(
        self, url: str, novel_name: str, save_per_chapter: bool
    ) -> List[ChapterData]: ...

    """
    Starting from the last stored chapter, keep following the 'nextchap' link
    until the link has the 'isDisabled' class. For each chapter page,
    fetch and save a ChapterData object.
    """

    def parse_chapters_with_notify(
        self, url: str, novel_name: str, save_per_chapter: bool
    ) -> Generator[str, None, List[ChapterData]]: ...

    def update_novel(
        self, novel_name: str, novel_url: str, last_chapter_url: str
    ) -> Tuple[str, List[ChapterData]]: ...

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
