from scrapper import config
from urllib.parse import urlparse
from scrapper.helpers import utils
from scrapper.modules.factories.factory import get_parser, SkipDuplicate
from scrapper.cache.db_cache import NovelDataCache
import argparse


def single_scrapper():
    parser = argparse.ArgumentParser(
        prog="scrapper", description="Single Novel scraper CLI"
    )
    parser.add_argument("url", help="URL of the novel page to start scraping at")

    args = parser.parse_args()

    novel_url = args.url
    cache = NovelDataCache(config.CACHE_DB_PATH)
    try:
        parser = get_parser(
            args.url,
            cache,
            skip_duplicates=SkipDuplicate.NOVEL,
            max_chapters_number=99999,
        )
        detail_tree = utils.fetch_page(novel_url)  # type: ignore
        novel_data = parser.parse_novel(detail_tree, novel_url)
        _ = parser.parse_chapters(novel_url, novel_data.title, True)

        print("✅ Finished. All novels saved.")

    except Exception as e:
        raise e


if __name__ == "__main__":
    single_scrapper()
