from scrapper.grpc import scrapper_pb2, scrapper_pb2_grpc
from concurrent import futures
import logging
import grpc
from typing import List
from scrapper import config
from urllib.parse import urlparse
from scrapper.helpers import utils
from scrapper.modules.factories.factory import get_parser, SkipDuplicate
from scrapper.cache.db_cache import NovelDataCache


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
                notify_buffer=notify_buffer,
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
    server.add_insecure_port("[::]:" + port)
    server.start()
    print("Server started, listening on " + port)
    server.wait_for_termination()


if __name__ == "__main__":
    logging.basicConfig()
    serve()
