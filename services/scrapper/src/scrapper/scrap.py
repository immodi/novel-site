from scrapper import config
from urllib.parse import urlparse
from scrapper.helpers import utils
from scrapper.modules.factories.factory import get_parser, SkipDuplicate
from scrapper.cache.db_cache import NovelDataCache
import argparse


def scrapper():
    parser = argparse.ArgumentParser(prog="scrapper", description="Novel scraper CLI")
    parser.add_argument("url", help="URL of the novel page to start scraping at")
    parser.add_argument(
        "start_page_num", type=int, help="Start page to scrape", default=1
    )
    parser.add_argument(
        "total_pages_num",
        type=int,
        help="Total amount of pages to scrape",
        default=1,
    )
    parser.add_argument(
        "max_novel_chapters_num",
        type=int,
        help="Max chapters number each novel is allowed to have",
        default=200,
    )

    args = parser.parse_args()
    list_urls = [
        f"{args.url}?page={i}"
        for i in range(args.start_page_num, args.start_page_num + args.total_pages_num)
    ]

    parsed = urlparse(args.url)
    base_url = f"{parsed.scheme}://{parsed.netloc}"

    config.BASE_URL = base_url
    cache = NovelDataCache(config.CACHE_DB_PATH)
    try:
        print("📥 Fetching novel list...")
        list_tree = utils.fetch_page(list_urls)  # type: ignore

        parser = get_parser(
            args.url,
            cache,
            SkipDuplicate.NOVEL,
            max_chapters_number=args.max_novel_chapters_num,
        )
        novels = parser.parse_list_of_novels(list_tree)
        for novel in novels:
            print(f"📖 Fetching {novel.title} → {novel.url}")
            detail_tree = utils.fetch_page(novel.url)  # type: ignore
            _ = parser.parse_novel(detail_tree, novel.url)
            _ = parser.parse_chapters(novel.url, novel.title, True)

        print("✅ Finished. All novels saved.")

    except Exception as e:
        raise e


if __name__ == "__main__":
    scrapper()
