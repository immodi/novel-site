from botasaurus.request import request, Request 
from botasaurus.soupify import soupify
from botasaurus.browser import browser, Driver
from urllib.parse import urljoin
from scrapper.helpers import utils, saver, helpers
from scrapper.modules.factories.factory import Parser, SkipDuplicate
from scrapper.helpers.utils import safe_text
from scrapper import config
from typing import List, Union
from scrapper.datatypes.novel import ChapterData, NovelData, NovelLink
from lxml import html
from scrapper.modules.utils.dataclasses.novel_fire_chapter import NovelFireChapter


class NovelFireParser(Parser):
    def __init__(self,  max_chapters_number: int, skip_duplicates: SkipDuplicate = SkipDuplicate.NONE):
        self.max_chapters_number = max_chapters_number
        self.skip_duplicates = skip_duplicates

    def parse_list_of_novels(
        self, tree: Union[html.HtmlElement, List[html.HtmlElement]]
    ) -> List[NovelLink]:
        novels = []

        # Normalize into a list
        trees = tree if isinstance(tree, list) else [tree]

        for t in trees:
            for novel_el in t.cssselect(".novel-list > *"):
                title_el = novel_el.cssselect("h4")
                link_el = novel_el.cssselect("a")

                if not title_el or not link_el:
                    continue

                title = title_el[0].text_content().strip()
                url = urljoin(config.BASE_URL, link_el[0].get("href"))

                if self.skip_duplicates == SkipDuplicate.NOVEL and self.novel_exists(title):
                    print(f"Novel {title} already exists. Skipping...")
                    continue

                novels.append({"title": title, "url": url})

        return novels

    def parse_novel(
        self, tree: html.HtmlElement, url: str, save_image: bool = True
    ) -> NovelData:
        # Title (safe_text in case of strange formatting)
        title_el = tree.cssselect(".novel-title")
        title = safe_text(title_el[0]) if len(title_el) > 0 else ""

        # Author
        author_el = tree.cssselect(".author a")
        author = author_el[0].text_content().strip() if len(author_el) > 0 else ""

        # Status
        status_el = tree.cssselect(".header-stats")
        if len(status_el) > 0 and len(status_el[0].getchildren()) >= 4:
            strong_el = status_el[0].getchildren()[3].cssselect("strong")
            status = strong_el[0].text_content().strip() if len(strong_el) > 0 else ""
        else:
            status = ""

        # Genres
        genres_el = tree.cssselect(".categories ul li")
        genres = [g.text_content().strip() for g in genres_el] if len(genres_el) > 0 else []

        # Tags
        tags_el = tree.cssselect("ul.content li")
        tags = [t.text_content().strip() for t in tags_el] if len(tags_el) > 0 else []

        # Cover image
        cover_el = tree.cssselect(".glass-background img")
        cover_image = cover_el[0].get("src") if len(cover_el) > 0 else ""

        # Description
        desc_el = tree.cssselect(".content")
        if len(desc_el) > 0:
            description = desc_el[0].text_content().split("Show More")[0].strip()
        else:
            description = ""

        novel_data: NovelData = {
            "title": title or "",
            "author": author,
            "genres": genres,
            "status": status,
            "tags": tags,
            "cover_image": cover_image,
            "description": description,
            "url": url,
        }

        saver.save_item(novel_data, f"{config.OUTPUT_DIR}/novels")
        if save_image and cover_image:
            helpers.download_image(
                cover_image, f"{config.OUTPUT_DIR}/covers", title or "cover"
            )

        return novel_data

    def parse_chapters(
        self, url: str, novel_name: str, save_per_chapter: bool
    ) -> List[ChapterData]:
        ajax_url = f"{url}/chapters"
        tree: html.HtmlElement = utils.fetch_page(ajax_url)  # type: ignore

        chapters: List[ChapterData] = []

        # Find the last chapter number
        header_el = tree.cssselect("header.container")
        if not header_el:
            return chapters

        last_link = header_el[0].getchildren()[-1].cssselect("a")
        if not last_link:
            return chapters

        href = last_link[0].get("href")
        if not href or "chapter-" not in href:
            return chapters

        base_url, last_chap = href.split("chapter-")
        try:
            last_chapter_number = int(last_chap)
        except ValueError:
            return chapters

        if last_chapter_number > self.max_chapters_number:
            print(
                f"Maximum number of chapters ({self.max_chapters_number}) exceeded. Skipping..."
            )
            self.clean_up_novel(novel_name)
            return chapters


        # Build chapter links (ascending order)
        for i in range(1, last_chapter_number + 1):
            chapter_link = f"{base_url}chapter-{i}"

            chapter_data = scrape_chapter_with_request(chapter_link)  # type: ignore
            # chapter_data = scrape_chapter(chapter_link)  # type: ignore
            
            if self.skip_duplicates == SkipDuplicate.CHAPTER and self.chapter_exists(
                chapter_data.title, utils.slugify(novel_name)
                
            ):
                print(f"Chapter {chapter_data.title} already exists. Skipping...")
                continue

            chapter: ChapterData = {
                "title": chapter_data.title,
                "content": chapter_data.content,
                "url": chapter_link,
            }

            chapters.append(chapter)

            if save_per_chapter:
                saver.save_item(
                    chapter, f"{config.OUTPUT_DIR}/chapters/{utils.slugify(novel_name)}"
                )

            print(
                f"--> Fetched {chapter_data.title} of length {utils.bold_green(len(chapter_data.content))} from url {chapter_link}."
            )

        if not save_per_chapter:
            for chapter in chapters:
                saver.save_item(
                    chapter, f"{config.OUTPUT_DIR}/chapters/{utils.slugify(novel_name)}"
                )

        return chapters


@browser(output=None, headless=False, max_retry=10)
def scrape_chapter(driver: Driver, url: str) -> NovelFireChapter:
    driver.google_get(url)
    title = driver.get_text(".chapter-title")
    content = driver.get_text("#content")
    return NovelFireChapter(title=title, content=content)



@request(output=None, max_retry=10)
def scrape_chapter_with_request(request: Request, url: str) -> NovelFireChapter:
    response = request.get(url, timeout=30)
    soup = soupify(response)

    if soup.get_text() == "Just a moment...Enable JavaScript and cookies to continue":
        return scrape_chapter(url)  # type: ignore

    chapter_title = soup.find("span", class_="chapter-title")
    title = chapter_title.get_text().strip() if chapter_title else ""
    
    chapter_content = soup.find(id="content")
    content = chapter_content.get_text().strip() if chapter_content else ""

    return NovelFireChapter(title=title, content=content)
