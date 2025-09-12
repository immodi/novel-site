from scrapper import config
from scrapper.modules.updater.list_novel_chapters_dirs import list_novel_dirs
from scrapper.modules.updater.collect_last_chapter_urls import collect_last_chapter_urls


def update():

    novel_dirs = list_novel_dirs(f"{config.OUTPUT_DIR}/chapters")
    last_chapters = collect_last_chapter_urls(novel_dirs)

    for item in last_chapters:
        print(item.novel_name, "->", item.last_chapter_url)


if __name__ == "__main__":
    update()
