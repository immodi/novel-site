from dataclasses import dataclass
from typing import Optional


@dataclass
class GetLastChapterByIdRequest:
    novel_id: int


@dataclass
class GetLastChapterByNameRequest:
    name: str


@dataclass
class GetLastChapterResponse:
    success: bool
    message: Optional[str] = None
    novel_id: Optional[int] = None
    last_chapter_number: Optional[int] = None
    last_chapter_name: Optional[str] = None
