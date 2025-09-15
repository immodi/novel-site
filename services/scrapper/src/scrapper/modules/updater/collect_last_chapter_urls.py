from pathlib import Path
from typing import List
from pathlib import Path
from scrapper.datatypes.last_chapter_url_data import NovelLastChapter
from scrapper.modules.updater.get_last_chapter_url import get_last_chapter_data


def collect_last_chapter_urls(novel_dirs: List[Path]) -> List[NovelLastChapter]:
    """
    Given a list of novel directories, find the last chapter URL for each
    and return a list of NovelLastChapter objects.
    """
    results: List[NovelLastChapter] = []

    for novel_dir in novel_dirs:
        last_chapter_data = get_last_chapter_data(str(novel_dir))
        if last_chapter_data is None:
            print(f"Skipping {novel_dir.name}: no valid chapters found.")
        else:
            last_url, last_title = last_chapter_data
            results.append(
                NovelLastChapter(
                    novel_name=novel_dir.name,
                    last_chapter_url=last_url,
                    last_chapter_name=last_title,
                )
            )
    return results
