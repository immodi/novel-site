from dataclasses import dataclass
from typing import Optional


@dataclass
class LoadNovelResponse:
    success: bool
    message: Optional[str] = None
    novel_id: Optional[int] = None
