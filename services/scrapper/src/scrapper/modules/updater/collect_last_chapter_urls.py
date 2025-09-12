from pathlib import Path
from typing import List
from pathlib import Path
from scrapper.datatypes.last_chapter_url_data import NovelLastChapter
from scrapper.modules.updater.get_last_chapter_url import get_last_chapter_url


def collect_last_chapter_urls(novel_dirs: List[Path]) -> List[NovelLastChapter]:
    """
    Given a list of novel directories, find the last chapter URL for each
    and return a list of NovelLastChapter objects.
    """
    results: List[NovelLastChapter] = []

    for novel_dir in novel_dirs:
        last_url = get_last_chapter_url(str(novel_dir))
        if last_url:
            results.append(
                NovelLastChapter(novel_name=novel_dir.name, last_chapter_url=last_url)
            )
        else:
            print(f"Skipping {novel_dir.name}: no valid chapters found.")
    return results
