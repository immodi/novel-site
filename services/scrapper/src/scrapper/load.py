import time
from pathlib import Path
from scrapper import config
from scrapper.datatypes.novel import NovelData
from scrapper.modules.loaders.load_from_json import load_from_json
from scrapper.modules.loaders.send_chapters_to_server import send_chapters_to_server
from scrapper.modules.loaders.send_to_server import send_novel_to_server
from scrapper.helpers.helpers import combine_json_objects_to_array


def loader(retry_delay: float = 2.0, max_retries: int = 5):
    novels_dir = Path(config.OUTPUT_DIR) / "novels"
    for json_file in novels_dir.glob("*.json"):
        novel = load_from_json(str(json_file), NovelData)

        retries = 0
        while True:
            resp = send_novel_to_server(novel)

            if resp.success:
                if resp.novel_id is None:
                    raise ValueError(
                        "Server response missing novel_id even though success=True"
                    )
                print(f"{json_file.name}: {resp}")
                novel_name = json_file.stem
                path, cleanup_callback = combine_json_objects_to_array(
                    f"{config.OUTPUT_DIR}/chapters/{novel_name}",
                )
                ch_resp = send_chapters_to_server(path, resp.novel_id)
                if ch_resp.success:
                    print(f"{json_file.name}: {ch_resp}")
                    cleanup_callback()
                    break
                break

            # skip if novel already exists
            if resp.message and "already exists" in resp.message.lower():
                print(f"{json_file.name}: already exists, skipping")
                break

            # retry if db is locked
            if resp.message and "database is locked" in resp.message.lower():
                retries += 1
                if retries > max_retries:
                    print(
                        f"{json_file.name}: database still locked after {max_retries} retries, skipping"
                    )
                    break
                print(
                    f"{json_file.name}: database locked, retrying in {retry_delay} seconds..."
                )
                time.sleep(retry_delay)
                continue

            # other errors
            print(f"{json_file.name}: failed to load novel: {resp}")
            break


if __name__ == "__main__":
    loader()
