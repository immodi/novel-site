import requests
from dataclasses import asdict
from scrapper.datatypes.novel import NovelData
from scrapper.datatypes.load_server_responses import LoadNovelResponse


def send_novel_to_server(
    novel: NovelData,
    image_path: str,
    url: str = "http://localhost:3000/load/novel",
    timeout: int = 10,
) -> LoadNovelResponse:
    try:
        payload = asdict(novel)

        # Convert payload dict into form-data fields
        data = {k: v for k, v in payload.items() if k != "cover_image"}

        # Open image file
        with open(f"{image_path}.jpg", "rb") as f:
            files = {"cover_image": (image_path.split("/")[-1], f, "image/jpeg")}
            response = requests.post(url, data=data, files=files, timeout=timeout)

        data = response.json()
        return LoadNovelResponse(
            success=data.get("success", False),
            message=data.get("message"),
            novel_id=data.get("novel_id"),
        )

    except requests.RequestException as e:
        return LoadNovelResponse(
            success=False, message=f"Request failed: {e}", novel_id=None
        )
    except Exception as e:
        return LoadNovelResponse(
            success=False, message=f"Unexpected error: {e}", novel_id=None
        )
