from typing import List
from dataclasses import dataclass


@dataclass
class NovelLink:
    title: str
    url: str


@dataclass
class NovelData:
    title: str
    author: str
    genres: List[str]
    status: str
    tags: List[str]
    cover_image: str
    description: str
    url: str


@dataclass
class ChapterData:
    title: str
    content: str
    url: str
