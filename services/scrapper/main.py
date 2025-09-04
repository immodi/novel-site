from modules import utils, saver, config
from modules.parsers.factory import get_parser


def main():
    print("📥 Fetching novel list...")
    list_tree = utils.fetch_page(config.LIST_URL)
    parser = get_parser(config.LIST_URL)
    novels = parser.parse_list_of_novels(list_tree)[:1]

    for novel in novels:
        print(f"📖 Fetching {novel['title']} → {novel['url']}")
        detail_tree = utils.fetch_page(novel["url"])
        novel_data = parser.parse_novel(detail_tree, novel["url"])
        chapters_data = parser.parse_chapters(novel["url"])
        saver.save_item(novel_data)
        for chapter in chapters_data:
            saver.save_item(chapter, f"chapters_json/{novel['title']}")

    print("✅ Finished. All novels saved.")


if __name__ == "__main__":
    main()
