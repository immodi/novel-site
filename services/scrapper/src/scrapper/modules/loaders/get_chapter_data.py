import json
from pathlib import Path
from scrapper.datatypes.last_chapter_url_data import NovelLastChapter


def get_chapter_data(
    novel_name: str,
    chapter_file: str,
    chapter_dir: str,
    ## dont actually need novel_name
) -> NovelLastChapter:
    """
    Parse a chapter JSON file and return a NovelLastChapter object.

    Args:
        chapter_file: JSON file name of the chapter
        chapter_dir: Directory where chapter JSON files are stored

    Returns:
        NovelLastChapter with last_chapter_url and last_chapter_name
    """
    chapter_dir_path = Path(chapter_dir)
    json_file = chapter_dir_path / novel_name / chapter_file

    if not json_file.exists():
        raise FileNotFoundError(
            f"No JSON file named {chapter_file} in {chapter_dir_path}"
        )

    with json_file.open("r", encoding="utf-8") as f:
        data = json.load(f)

    title = data.get("title")
    url = data.get("url")

    if not title or not url:
        raise ValueError(f"JSON file {chapter_file} does not contain 'title' or 'url'")

    return NovelLastChapter(
        novel_name=novel_name, last_chapter_url=url, last_chapter_name=title
    )
