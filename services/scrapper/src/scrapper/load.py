import time
from pathlib import Path
from scrapper import config
from scrapper.datatypes.novel import NovelData
from scrapper.modules.loaders.load_from_json import load_from_json
from scrapper.modules.loaders.send_to_server import send_novel_to_server


def loader(retry_delay: float = 2.0, max_retries: int = 5):
    novels_dir = Path(config.OUTPUT_DIR) / "novels"

    for json_file in novels_dir.glob("*.json"):
        novel = load_from_json(str(json_file), NovelData)

        retries = 0
        while True:
            resp = send_novel_to_server(novel)

            if resp.success:
                print(f"{json_file.name}: {resp}")
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
