import argparse
from time import sleep
from scrapper import config
from scrapper.helpers.move_files_up_and_delete import move_files_up_and_delete
from scrapper.helpers.combine_json_objects_into_array import (
    append__combine_json_objects_to_array,
)
from scrapper.modules.loaders.log_if_mismatch import log_if_mismatch
from scrapper.modules.loaders.get_last_chapter import LastChapterClient
from scrapper.modules.loaders.send_chapters_to_server import append_chapters_to_server
from scrapper.modules.factories.factory import get_parser, SkipDuplicate
from scrapper.cache.db_cache import NovelDataCache


def update():
    # Argument parsing inside the function
    parser = argparse.ArgumentParser(prog="updater", description="Novel updater CLI")
    parser.add_argument(
        "interval_hours",
        type=float,
        default=12,
        help="Interval in hours between updates (default: 12)",
    )
    args = parser.parse_args()
    interval_hours = args.interval_hours

    print(f"Updater started. Interval hours: {interval_hours}")

    cache = NovelDataCache(config.CACHE_DB_PATH)
    while True:
        try:
            print("Starting update...")

            novels = cache.get_all_novels()
            last_chapters = cache.get_last_chapters(novels)

            client = LastChapterClient()

            for chapter in last_chapters:
                db_last_chapter = client.get_by_name(
                    chapter.novel_name
                )  # get last chapter from db

                last_chapter_number = db_last_chapter.last_chapter_number
                novel_id = db_last_chapter.novel_id

                is_mismatched = log_if_mismatch(
                    remote=db_last_chapter,
                    local_name=chapter.last_chapter_name,
                    # for logging
                    novel_identifier=chapter.novel_name,
                )

                if is_mismatched:
                    print(f"Skipping {chapter.novel_name} due to mismatch")
                    continue

                if last_chapter_number is None or novel_id is None:
                    print(f"Skipping {chapter.novel_name} due to missing data")
                    continue

                parser_obj = get_parser(
                    chapter.last_chapter_url,
                    cache,
                    SkipDuplicate.NONE,  # doesnt matter
                )

                novel_url = cache.get_novel_url_by_chapter(chapter.last_chapter_url)
                if novel_url is None:
                    print(f"Skipping {chapter.novel_name} due to missing novel URL")
                    continue

                save_dir, new_chapters_list = parser_obj.update_novel(
                    chapter.novel_name, novel_url, chapter.last_chapter_url
                )
                if len(new_chapters_list) == 0:
                    continue

                numbers = [
                    i
                    for i in range(
                        last_chapter_number + 1,
                        last_chapter_number + len(new_chapters_list) + 1,
                    )
                ]

                path, cleanup_callback = append__combine_json_objects_to_array(
                    save_dir, numbers
                )
                ch_resp = append_chapters_to_server(path, novel_id)
                if ch_resp.success:
                    print(f"{chapter.novel_name}: {ch_resp}")
                    cleanup_callback()
                    move_files_up_and_delete(save_dir)
                else:
                    print(f"Failed to update chapters for novel {chapter.novel_name}")

            print("Update complete.")
        except Exception as e:
            print(f"Error occurred during update: {e}")

        print(f"Sleeping for {interval_hours} hours...")
        sleep(interval_hours * 60 * 60)


if __name__ == "__main__":
    update()
