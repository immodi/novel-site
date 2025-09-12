from typing import List
from pathlib import Path


def list_novel_dirs(base_chapters_dir: str) -> List[Path]:
    """
    Return all subdirectories inside the chapters directory.
    Each subdirectory represents a single novel.
    """
    base_path = Path(base_chapters_dir)
    if not base_path.exists():
        print(f"Directory not found: {base_chapters_dir}")
        return []

    return [d for d in base_path.iterdir() if d.is_dir()]
