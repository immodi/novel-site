from scrapper import config
from scrapper.helpers import utils, helpers
from scrapper.modules.parsers.factory import get_parser
from scrapper.modules.parsers.interface import SkipDuplicate


def main():
    print("📥 Fetching novel list...")
    list_tree = utils.fetch_page(config.LIST_URL)  # type: ignore
    parser = get_parser(config.LIST_URL, SkipDuplicate.CHAPTER)
    novels = parser.parse_list_of_novels(list_tree)
    for novel in novels:
        print(f"📖 Fetching {novel['title']} → {novel['url']}")
        detail_tree = utils.fetch_page(novel["url"])  # type: ignore
        _ = parser.parse_novel(detail_tree, novel["url"])
        _ = parser.parse_chapters(novel["url"], novel["title"], True)

    print("✅ Finished. All novels saved.")


if __name__ == "__main__":
    main()
