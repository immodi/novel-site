from typing import TypedDict, List


class NovelLink(TypedDict):
    title: str
    url: str


class NovelData(TypedDict):
    title: str
    author: str
    genres: List[str]
    status: str
    tags: List[str]
    cover_image: str
    description: str
    url: str


class ChapterData(TypedDict):
    title: str
    content: str
    url: str
