from typing import TypedDict, List, Optional


class NovelLink(TypedDict):
    title: str
    url: str


class NovelData(TypedDict, total=False):
    title: Optional[str]
    author: Optional[str]
    genres: List[str]
    status: Optional[str]
    tags: List[str]
    cover_image: Optional[str]
    description: Optional[str]
    url: str


class ChapterData(TypedDict):
    title: str
    content: str
    url: str
