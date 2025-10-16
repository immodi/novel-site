import requests
import json
from pathlib import Path
from scrapper.datatypes.load_server_responses import LoadNovelResponse


def send_chapters_to_server(
    file_path: str, novel_id: int, url: str = "http://localhost:3000/load/chapters"
) -> LoadNovelResponse:
    path_obj = Path(file_path)
    if not path_obj.exists() or not path_obj.is_file():
        raise FileNotFoundError(f"File not found: {path_obj}")

    metadata = {"novel_id": novel_id}

    with open(path_obj, "rb") as f:
        files = {"file": (path_obj.name, f, "application/json")}
        data = {"metadata": json.dumps(metadata)}

        response = requests.post(url, files=files, data=data)
        response.raise_for_status()
        res_json = response.json()

        # Convert response JSON to dataclass
        return LoadNovelResponse(
            success=res_json.get("success", False),
            message=res_json.get("message", ""),
            novel_id=res_json.get("novel_id", 0),
        )


def append_chapters_to_server(
    file_path: str,
    novel_id: int,
    url: str = "http://localhost:3000/load/append-chapters",
) -> LoadNovelResponse:
    """
    Send a JSON file of chapters to the Go endpoint /load/append-chapters.

    Parameters
    ----------
    file_path : str
        Path to the JSON file containing a list of chapters.
    novel_id : int
        ID of the novel to which the chapters belong.
    url : str, optional
        Full endpoint URL. Defaults to http://localhost:3000/load/append-chapters.

    Returns
    -------
    LoadNovelResponse
        Dataclass representing the API response.
    """
    path_obj = Path(file_path)
    if not path_obj.exists() or not path_obj.is_file():
        raise FileNotFoundError(f"File not found: {path_obj}")

    metadata = {"novel_id": novel_id}

    with open(path_obj, "rb") as f:
        files = {"file": (path_obj.name, f, "application/json")}
        data = {"metadata": json.dumps(metadata)}

        response = requests.post(url, files=files, data=data)
        # response.raise_for_status()
        res_json = response.json()

        return LoadNovelResponse(
            success=res_json.get("success", False),
            message=res_json.get("message", ""),
            novel_id=res_json.get("novel_id", 0),
        )
