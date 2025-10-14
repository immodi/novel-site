from scrapper import config
from scrapper.modules.factories.factory import get_parser, SkipDuplicate
from scrapper.cache.db_cache import NovelDataCache
from scrapper.helpers.move_files_up_and_delete import move_files_up_and_delete
from scrapper.helpers.combine_json_objects_into_array import (
    append__combine_json_objects_to_array,
)
from scrapper.modules.loaders.log_if_mismatch import log_if_mismatch
from scrapper.modules.loaders.get_last_chapter import LastChapterClient
from scrapper.modules.loaders.send_chapters_to_server import append_chapters_to_server


def run_update(broadcast):
    cache = NovelDataCache(config.CACHE_DB_PATH)
    broadcast("Starting update...")

    novels = cache.get_all_novels()
    last_chapters = cache.get_last_chapters(novels)
    client = LastChapterClient()

    for chapter in last_chapters:
        db_last_chapter = client.get_by_name(chapter.novel_name)
        if not db_last_chapter.success:
            broadcast(f"Skipping {chapter.novel_name} due to invalid novel name")
            continue

        last_chapter_number = db_last_chapter.last_chapter_number
        novel_id = db_last_chapter.novel_id

        is_mismatched = log_if_mismatch(
            remote=db_last_chapter,
            local_name=chapter.last_chapter_name,
            novel_identifier=chapter.novel_name,
        )

        if is_mismatched:
            broadcast(f"Skipping {chapter.novel_name} due to mismatch")
            continue

        if last_chapter_number is None or novel_id is None:
            msg = f"Skipping {chapter.novel_name} due to missing data"
            broadcast(msg)
            continue

        parser_obj = get_parser(chapter.last_chapter_url, cache, SkipDuplicate.NONE)
        novel_url = cache.get_novel_url_by_chapter(chapter.last_chapter_url)

        if novel_url is None:
            broadcast(f"Skipping {chapter.novel_name} due to missing novel URL")
            continue

        chapter_gen = parser_obj.update_novel_with_notify(
            chapter.novel_name, novel_url, chapter.last_chapter_url
        )

        try:
            for msg in chapter_gen:
                broadcast(msg)
        except Exception as e:
            broadcast(f"Error during update: {e}")
            continue

        try:
            result = chapter_gen.send(None)
        except StopIteration as stop:
            result = stop.value

        if not result:
            continue

        save_dir, new_chapters_list = result
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
            broadcast(f"{chapter.novel_name}: {ch_resp}")
            cleanup_callback()
            move_files_up_and_delete(save_dir)
        else:
            broadcast(f"Failed to update chapters for novel {chapter.novel_name}")

    broadcast("Update complete.")
