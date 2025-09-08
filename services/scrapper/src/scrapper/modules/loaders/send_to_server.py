import requests
from dataclasses import asdict
from scrapper.datatypes.novel import NovelData
from scrapper.datatypes.load_server_responses import LoadNovelResponse


def send_novel_to_server(
    novel: NovelData, url: str = "http://localhost:3000/load-novel", timeout: int = 10
) -> LoadNovelResponse:
    try:
        payload = asdict(novel)
        response = requests.post(url, json=payload, timeout=timeout)
        data = response.json()
        return LoadNovelResponse(
            success=data.get("success", False),
            message=data.get("message"),
            novel_id=data.get("novel_id"),
        )
    except requests.RequestException as e:
        # always return a LoadNovelResponse object
        return LoadNovelResponse(
            success=False, message=f"Request failed: {e}", novel_id=None
        )
