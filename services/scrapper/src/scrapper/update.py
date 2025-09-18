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
from scrapper.modules.updater.list_novel_chapters_dirs import list_novel_dirs
from scrapper.modules.updater.collect_last_chapter_urls import collect_last_chapter_urls
from scrapper.modules.factories.factory import get_parser, SkipDuplicate


def update():
    novel_dirs = list_novel_dirs(f"{config.OUTPUT_DIR}/chapters")
    last_chapters = collect_last_chapter_urls(novel_dirs)
    client = LastChapterClient()

    for chapter in last_chapters:
        db_last_chapter = client.get_by_name(chapter.novel_name)
        last_chapter_number = db_last_chapter.last_chapter_number
        novel_id = db_last_chapter.novel_id

        is_mismatched = log_if_mismatch(
            remote=db_last_chapter,
            local_name=chapter.last_chapter_name,
            novel_identifier=chapter.novel_name,
        )

        # Skip if there's a mismatch AND either value is missing
        if is_mismatched:
            continue

        if last_chapter_number is None or novel_id is None:
            continue

        parser = get_parser(
            chapter.last_chapter_url,
            SkipDuplicate.CHAPTER,
        )
        save_dir, new_chapters_list = parser.update_novel(
            chapter.novel_name, chapter.last_chapter_url
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


if __name__ == "__main__":
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

    while True:
        try:
            print("Starting update...")
            update()
            print("Update complete.")
        except Exception as e:
            print(f"Error occurred during update: {e}")

        print(f"Sleeping for {interval_hours} hours...")
        sleep(interval_hours * 60 * 60)
