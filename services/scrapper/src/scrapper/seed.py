from pathlib import Path
from scrapper import config
from scrapper.datatypes.novel import NovelData
from scrapper.modules.loaders.load_from_json import load_from_json
from scrapper.cache.db_cache import NovelDataCache
from scrapper.modules.updater.get_novel_url_by_dir_name import get_novel_url
from scrapper.modules.updater.list_novel_chapters_dirs import list_novel_dirs
from scrapper.modules.updater.collect_last_chapter_urls import collect_last_chapter_urls


def seed():
    cache = NovelDataCache(config.CACHE_DB_PATH)
    novels_dir = Path(config.OUTPUT_DIR) / "novels"

    print("Saving novels...")
    for json_file in novels_dir.glob("*.json"):
        novel = load_from_json(str(json_file), NovelData)
        print(f"Saving novel {novel.title}")
        cache.save_novel(novel)

    print("Saving last chapters...")
    novel_dirs = list_novel_dirs(f"{config.OUTPUT_DIR}/chapters")
    last_chapters = collect_last_chapter_urls(novel_dirs)

    for chapter in last_chapters:
        novel_url = get_novel_url(f"{chapter.novel_name}.json")
        if novel_url is None:
            print(f"Novel not found for chapter {chapter.last_chapter_url}")
            continue

        print(
            f"Saving last chapter {chapter.last_chapter_name} for {chapter.novel_name}"
        )

        cache.save_last_chapter(
            novel_url=novel_url,
            chapter_url=chapter.last_chapter_url,
            chapter_name=chapter.last_chapter_name,
        )


if __name__ == "__main__":
    seed()
