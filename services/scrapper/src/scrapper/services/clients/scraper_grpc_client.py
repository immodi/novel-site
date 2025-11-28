from __future__ import print_function
import logging
import grpc
from scrapper.grpc_services import scrapper_pb2, scrapper_pb2_grpc


def run():
    with grpc.insecure_channel("localhost:50051") as channel:
        stub = scrapper_pb2_grpc.ScrapperServiceStub(channel)

        request = scrapper_pb2.ScrapeRequest(
            url="https://novelfire.net/genre-all/sort-popular/status-all/all-novel",
            start_page_num=20,
            total_pages_num=1,
            max_novel_chapters_num=200,
        )

        print("Starting ScrapeNovels stream...")
        try:
            for response in stub.ScrapeNovels(request):
                # Check which field is populated
                if response.novel_link.title or response.novel_link.url:
                    novel = response.novel_link
                    print(f"Received novel: {novel.title} → {novel.url}")
                elif response.chapter_message.message:
                    chapter = response.chapter_message
                    print(f"Chapter message: {chapter.message}")
                else:
                    print("Received empty ScrapeResponse")
        except grpc.RpcError as e:
            print(f"gRPC error: {e.code()} - {e.details()}")


if __name__ == "__main__":
    logging.basicConfig()
    run()
