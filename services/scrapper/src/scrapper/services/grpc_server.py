from scrapper.grpc import scrapper_pb2, scrapper_pb2_grpc, updater_pb2, updater_pb2_grpc
from concurrent import futures
import logging
import grpc
from typing import List
from scrapper import config
from urllib.parse import urlparse
from scrapper.helpers import utils
from scrapper.modules.factories.factory import get_parser, SkipDuplicate
from scrapper.cache.db_cache import NovelDataCache
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


class UpdaterService(updater_pb2_grpc.UpdaterServiceServicer):
    def UpdateNovels(self, request, context):
        try:
            interval_hours = request.interval_hours
            yield updater_pb2.UpdateResponse(
                message=f"Updater started. Interval hours: {interval_hours}"
            )

            cache = NovelDataCache(config.CACHE_DB_PATH)
            try:
                yield updater_pb2.UpdateResponse(message="Starting update...")
                novels = cache.get_all_novels()
                last_chapters = cache.get_last_chapters(novels)

                client = LastChapterClient()
                for chapter in last_chapters:
                    db_last_chapter = client.get_by_name(chapter.novel_name)

                    if not db_last_chapter.success:
                        msg = f"Skipping {chapter.novel_name} due to invalid novel name"
                        yield updater_pb2.UpdateResponse(message=msg)
                        continue

                    last_chapter_number = db_last_chapter.last_chapter_number
                    novel_id = db_last_chapter.novel_id

                    is_mismatched = log_if_mismatch(
                        remote=db_last_chapter,
                        local_name=chapter.last_chapter_name,
                        novel_identifier=chapter.novel_name,
                    )

                    if is_mismatched:
                        msg = f"Skipping {chapter.novel_name} due to mismatch"
                        yield updater_pb2.UpdateResponse(message=msg)
                        continue

                    if last_chapter_number is None or novel_id is None:
                        msg = f"Skipping {chapter.novel_name} due to missing data"
                        yield updater_pb2.UpdateResponse(message=msg)
                        continue

                    parser_obj = get_parser(
                        chapter.last_chapter_url,
                        cache,
                        SkipDuplicate.NONE,  # doesn’t matter
                    )

                    novel_url = cache.get_novel_url_by_chapter(chapter.last_chapter_url)
                    if novel_url is None:
                        msg = f"Skipping {chapter.novel_name} due to missing novel URL"
                        yield updater_pb2.UpdateResponse(message=msg)
                        continue

                    # Stream messages as the update runs
                    chapter_gen = parser_obj.update_novel_with_notify(
                        chapter.novel_name, novel_url, chapter.last_chapter_url
                    )

                    # Iterate over progress messages
                    try:
                        for msg in chapter_gen:
                            yield updater_pb2.UpdateResponse(message=msg)
                    except Exception as e:
                        yield updater_pb2.UpdateResponse(
                            message=f"Error during update: {e}"
                        )
                        continue

                    # Capture the return value (save_dir, new_chapters_list)
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
                        msg = f"{chapter.novel_name}: {ch_resp}"
                        yield updater_pb2.UpdateResponse(message=msg)
                        cleanup_callback()
                        move_files_up_and_delete(save_dir)
                    else:
                        msg = (
                            f"Failed to update chapters for novel {chapter.novel_name}"
                        )
                        yield updater_pb2.UpdateResponse(message=msg)

                msg = "Update complete."
                yield updater_pb2.UpdateResponse(message=msg)

            except Exception as e:
                msg = f"Error occurred during update: {e}"
                yield updater_pb2.UpdateResponse(message=msg)

            msg = "Update Done!"
            yield updater_pb2.UpdateResponse(message=msg)

        except Exception as e:
            context.set_details(str(e))
            context.set_code(grpc.StatusCode.INTERNAL)
            return


class ScrapperService(scrapper_pb2_grpc.ScrapperServiceServicer):
    def ScrapeNovels(self, request, context):
        notify_buffer: List[str] = []

        # Build list of page URLs
        list_urls = [
            f"{request.url}?page={i}"
            for i in range(
                request.start_page_num, request.start_page_num + request.total_pages_num
            )
        ]

        parsed = urlparse(request.url)
        base_url = f"{parsed.scheme}://{parsed.netloc}"

        config.BASE_URL = base_url
        cache = NovelDataCache(config.CACHE_DB_PATH)

        try:
            print("📥 Fetching novel list...")
            list_tree = utils.fetch_page(list_urls)  # type: ignore

            parser = get_parser(
                request.url,
                cache,
                SkipDuplicate.NOVEL,
                notify_buffer=notify_buffer,  # type: ignore
                max_chapters_number=request.max_novel_chapters_num,
            )
            novels = parser.parse_list_of_novels(list_tree)

            for novel in novels:
                print(f"📖 Fetching {novel.title} → {novel.url}")
                detail_tree = utils.fetch_page(novel.url)  # type: ignore
                _ = parser.parse_novel(detail_tree, novel.url)

                # Stream chapter messages while fetching
                chapter_gen = parser.parse_chapters_with_notify(
                    novel.url, novel.title, True
                )
                for msg in chapter_gen:
                    yield scrapper_pb2.ScrapeResponse(
                        novel_link=scrapper_pb2.NovelLink(title="", url=""),
                        chapter_message=scrapper_pb2.ChapterMessage(message=msg),
                    )

                # Finally, yield the novel itself
                yield scrapper_pb2.ScrapeResponse(
                    novel_link=scrapper_pb2.NovelLink(title=novel.title, url=novel.url),
                    chapter_message=scrapper_pb2.ChapterMessage(message=""),
                )

        except Exception as e:
            context.set_details(str(e))
            context.set_code(grpc.StatusCode.INTERNAL)
            return


def serve():
    port = "50051"
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    scrapper_pb2_grpc.add_ScrapperServiceServicer_to_server(ScrapperService(), server)
    updater_pb2_grpc.add_UpdaterServiceServicer_to_server(UpdaterService(), server)
    server.add_insecure_port("[::]:" + port)
    server.start()
    print("Server started, listening on " + port)
    server.wait_for_termination()


if __name__ == "__main__":
    logging.basicConfig()
    serve()
