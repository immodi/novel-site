import json
from pathlib import Path
from scrapper.datatypes.novel import NovelData
from scrapper import config


def get_novel_url(
    file_name: str, novels_dir: str = f"{config.OUTPUT_DIR}/novels"
) -> str:
    """
    Read a JSON file in the novels directory and return the novel's URL.

    Args:
        file_name: Name of the JSON file (with or without .json extension)
        novels_dir: Directory where novel JSON files are stored
    """
    novels_dir_path = Path(novels_dir)
    file_name = file_name if file_name.endswith(".json") else f"{file_name}.json"
    json_file = novels_dir_path / file_name

    if not json_file.exists():
        raise FileNotFoundError(f"No JSON file named {file_name} in {novels_dir_path}")

    with json_file.open("r", encoding="utf-8") as f:
        data = json.load(f)

    novel = NovelData(
        title=data.get("title", ""),
        author=data.get("author", ""),
        genres=data.get("genres", []),
        status=data.get("status", ""),
        tags=data.get("tags", []),
        cover_image=data.get("cover_image", ""),
        description=data.get("description", ""),
        url=data.get("url", ""),
    )

    return novel.url
