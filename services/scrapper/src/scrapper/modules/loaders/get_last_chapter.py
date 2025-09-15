import json
import requests
from pathlib import Path

from scrapper.datatypes.get_last_chapter import (
    GetLastChapterByIdRequest,
    GetLastChapterByNameRequest,
    GetLastChapterResponse,
)
from scrapper.config import OUTPUT_DIR


class LastChapterClient:
    def __init__(self, base_url: str = "http://localhost:3000"):
        """
        Client for the /load/last-chapter endpoints.

        Parameters
        ----------
        base_url : str
            Root URL of your Go server. Defaults to http://localhost:3000
        """
        self.base_url = base_url.rstrip("/")

    # ------------------------------------------------------------------ #
    # Helper to send POST and print server JSON if status != 2xx
    # ------------------------------------------------------------------ #
    def _post_json(self, url: str, payload: dict) -> dict:
        r = requests.post(url, json=payload)
        if not r.ok:
            try:
                # print server JSON message for easier debugging
                print(f"[{r.status_code}] {url} ->", r.json())
            except Exception:
                print(f"[{r.status_code}] {url} ->", r.text)
            r.raise_for_status()
        return r.json()

    # ------------------------------------------------------------------ #
    def get_by_id(self, novel_id: int) -> GetLastChapterResponse:
        """POST /load/last-chapter/id"""
        url = f"{self.base_url}/load/last-chapter/id"
        payload = GetLastChapterByIdRequest(novel_id=novel_id)
        data = self._post_json(url, payload.__dict__)
        return GetLastChapterResponse(**data)

    def get_by_name(self, file_name: str) -> GetLastChapterResponse:
        """
        POST /load/last-chapter/name

        Parameters
        ----------
        file_name : str
            Name of the local novel JSON file (without .json extension).
            The file must contain a 'title' field whose value is the
            actual novel title stored in the database.
        """
        novel_json_path = Path(OUTPUT_DIR) / "novels" / f"{file_name}.json"

        with open(novel_json_path, "r", encoding="utf-8") as f:
            novel_data = json.load(f)

        # Extract the real title field from the JSON
        real_title = novel_data["title"]

        url = f"{self.base_url}/load/last-chapter/name"
        payload = GetLastChapterByNameRequest(name=real_title)
        data = self._post_json(url, payload.__dict__)
        return GetLastChapterResponse(**data)
