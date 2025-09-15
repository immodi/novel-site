import json
from pathlib import Path
from typing import Callable, Tuple
import tempfile
from scrapper.helpers.helpers import CHAPTER_REGEX


def append__combine_json_objects_to_array(
    dir_path: str, numbers: list[int]
) -> Tuple[str, Callable[[], None]]:
    """
    Combine all JSON files in a directory into a single JSON array file,
    and assign a 'number' field to each object from the supplied numbers list,
    in ascending order of chapter URL numbers (not necessarily matching).

    Returns:
        (temp_file_path, cleanup_callback)
    """
    dir_path_obj = Path(dir_path)
    json_entries: list[tuple[int, dict]] = []

    for json_file in dir_path_obj.glob("*.json"):
        try:
            with open(json_file, "r", encoding="utf-8") as f:
                data = json.load(f)
            if isinstance(data, dict) and "url" in data:
                match = CHAPTER_REGEX.search(data["url"])
                if match:
                    chapter_num = int(match.group(1))
                    json_entries.append((chapter_num, data))
                else:
                    print(
                        f"Skipping {json_file.name}: could not find chapter number in URL"
                    )
        except Exception as e:
            print(f"Skipping {json_file.name}: {e}")

    # Sort by the chapter number found in the URL
    json_entries.sort(key=lambda x: x[0])

    if len(numbers) < len(json_entries):
        raise ValueError(
            f"Not enough numbers supplied: have {len(numbers)}, need {len(json_entries)}"
        )

    combined: list[dict] = []
    for idx, (_, chapter_data) in enumerate(json_entries):
        # Attach the supplied number (by order)
        chapter_data["number"] = numbers[idx]
        combined.append(chapter_data)

    temp_dir = Path(tempfile.gettempdir())
    temp_file_path = temp_dir / f"{dir_path_obj.name}.json"

    with open(temp_file_path, "w", encoding="utf-8") as f:
        json.dump(combined, f, ensure_ascii=False, indent=2)

    print(
        f"Combined {len(combined)} JSON objects into {temp_file_path} "
        f"with supplied numbers."
    )

    def cleanup():
        if temp_file_path.exists():
            temp_file_path.unlink()
            print(f"Deleted temp file: {temp_file_path}")

    return str(temp_file_path), cleanup


def combine_json_objects_to_array(dir_path: str) -> Tuple[str, Callable[[], None]]:
    dir_path_obj = Path(dir_path)  # convert to Path
    combined = []

    # Collect JSON objects and pair with their chapter number
    json_entries = []
    for json_file in dir_path_obj.glob("*.json"):
        try:
            with open(json_file, "r", encoding="utf-8") as f:
                data = json.load(f)
                if isinstance(data, dict) and "url" in data:
                    match = CHAPTER_REGEX.search(data["url"])
                    if match:
                        chapter_num = int(match.group(1))
                        json_entries.append((chapter_num, data))
                    else:
                        print(
                            f"Skipping {json_file.name}: could not find chapter number in URL"
                        )
        except Exception as e:
            print(f"Skipping {json_file.name}: {e}")

    # Sort by chapter number
    json_entries.sort(key=lambda x: x[0])

    # Extract only the data part
    combined = [entry[1] for entry in json_entries]

    # Create a temp file in the system temp dir with the last dir name
    temp_dir = Path(tempfile.gettempdir())
    temp_file_path = temp_dir / f"{dir_path_obj.name}.json"

    with open(temp_file_path, "w", encoding="utf-8") as f:
        json.dump(combined, f)

    print(f"Combined {len(combined)} JSON objects into {temp_file_path}")

    # Cleanup callback
    def cleanup():
        if temp_file_path.exists():
            temp_file_path.unlink()
            print(f"Deleted temp file: {temp_file_path}")

    return str(temp_file_path), cleanup
