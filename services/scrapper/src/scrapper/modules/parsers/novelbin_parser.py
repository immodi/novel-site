from botasaurus.request import request, Request
from bs4 import BeautifulSoup, Tag
from botasaurus.soupify import soupify
from botasaurus.browser import browser, Driver
from urllib.parse import urljoin
from scrapper.helpers import utils, saver, helpers
from scrapper.modules.factories.factory import Parser, SkipDuplicate
from scrapper.helpers.utils import safe_text
from scrapper import config
from typing import List, Tuple, Union
from scrapper.datatypes.novel import ChapterData, NovelData, NovelLink
from lxml import html
from scrapper.cache.db_cache import NovelDataCache


class NovelBinParser(Parser):
    def __init__(
        self,
        max_chapters_number: int,
        cache: NovelDataCache,
        skip_duplicates: SkipDuplicate = SkipDuplicate.NONE,
    ):
        self.max_chapters_number = max_chapters_number
        self.cache = cache
        self.skip_duplicates = skip_duplicates

    def parse_list_of_novels(
        self, tree: Union[html.HtmlElement, List[html.HtmlElement]]
    ) -> List[NovelLink]:
        novels = []

        # Normalize into a list
        trees = tree if isinstance(tree, list) else [tree]

        for t in trees:
            for a in t.cssselect(".list-novel .row .novel-title > a"):
                title = a.text_content().strip()

                if self.skip_duplicates == SkipDuplicate.NOVEL and self.novel_exists(
                    title
                ):
                    print(f"Novel {title} already exists. Skipping...")
                    continue

                novels.append(
                    NovelLink(
                        title=title,
                        url=urljoin(config.BASE_URL, a.get("href")),
                    )
                )

        return novels

    def parse_novel(
        self, tree: html.HtmlElement, url: str, save_image: bool = True
    ) -> NovelData:
        """Extract data from a single novel page"""
        title = (
            safe_text(tree.cssselect(".title")[0]) if tree.cssselect(".title") else None
        )

        author, status = None, None
        genres: list[str] = []
        info_items = tree.cssselect(".info.info-meta li")

        for li in info_items:
            header = safe_text(li.cssselect("h3")[0]) if li.cssselect("h3") else ""
            header = header.lower().rstrip(":") if header else ""

            if header == "author":
                author = safe_text(li.cssselect("a")[0]) if li.cssselect("a") else None
            elif header == "genre":
                genres = [g for g in (safe_text(g) for g in li.cssselect("a")) if g]
            elif header == "status":
                status = safe_text(li.cssselect("a")[0]) if li.cssselect("a") else None

        tags: list[str] = [
            t for t in (safe_text(t) for t in tree.cssselect(".tag-container a")) if t
        ]

        img_el = tree.cssselect(".book > img")
        img = img_el[0].get("data-src") or img_el[0].get("src") if img_el else None

        description_url = url + "#tab-description-title"
        description_tree = utils.fetch_page(description_url)  # type: ignore
        desc_el = description_tree.cssselect(".desc-text")
        description = safe_text(desc_el[0]) if desc_el else None

        novel_data = NovelData(
            title=title or "",
            author=author or "",
            genres=genres or [],
            status=status or "",
            tags=tags or [],
            cover_image=img or "",
            description=description or "",
            url=url or "",
        )

        saver.save_item(novel_data, f"{config.OUTPUT_DIR}/novels")
        if save_image:
            helpers.download_image(
                novel_data.cover_image,
                f"{config.OUTPUT_DIR}/covers",
                novel_data.title,
            )

        self.cache.save_novel(novel_data)

        return novel_data

    def parse_chapters(
        self, url: str, novel_name: str, save_per_chapter: bool
    ) -> List[ChapterData]:
        """Fetch chapters using the first link and navigating via #next_chap button"""
        # Fetch chapter list via AJAX
        novel_id = url.rstrip("/").split("/")[-1]
        ajax_url = f"{config.BASE_URL}/ajax/chapter-archive?novelId={novel_id}"
        tree: html.HtmlElement = utils.fetch_page(ajax_url)  # type: ignore

        chapters: List[ChapterData] = []
        chapters_links = tree.cssselect(".panel-body a")

        if len(chapters_links) > self.max_chapters_number:
            print(
                f"Maximum number of chapters ({self.max_chapters_number}) exceeded. Skipping..."
            )
            self.clean_up_novel(novel_name)
            self.cache.remove_novel_by_url(url)
            return chapters

        for chapter_link in chapters_links:
            chapter_title = chapter_link.text_content().strip()
            if self.skip_duplicates == SkipDuplicate.CHAPTER and self.chapter_exists(
                chapter_title, utils.slugify(novel_name)
            ):
                print(f"Chapter {chapter_title} already exists. Skipping...")
                continue

            chapter_url = urljoin(config.BASE_URL, chapter_link.get("href"))
            chapter_content = scrape_chapter_with_request(chapter_url)  # type: ignore
            # chapter_content = scrape_chapter(chapter_url)  # type: ignore

            chapter = ChapterData(
                title=chapter_title,
                content=chapter_content,
                url=chapter_url,
            )
            chapters.append(chapter)

            if save_per_chapter:
                saver.save_item(
                    chapter, f"{config.OUTPUT_DIR}/chapters/{utils.slugify(novel_name)}"
                )
                self.cache.save_last_chapter(url, chapter.url, chapter.title)

            print(
                f"--> Fetched {chapter_title} of length {utils.bold_green(len(chapter_content))} from url {chapter_url}."
            )

        if not save_per_chapter:
            for chapter in chapters:
                saver.save_item(
                    chapter, f"{config.OUTPUT_DIR}/chapters/{utils.slugify(novel_name)}"
                )
            self.cache.save_last_chapter(url, chapters[-1].url, chapters[-1].title)

        return chapters

    def update_novel(
        self, novel_name: str, novel_url: str, last_chapter_url: str
    ) -> Tuple[str, List[ChapterData]]:
        return "", []


@browser(output=None, headless=True, max_retry=10)
def scrape_chapter(driver: Driver, url: str) -> str:
    driver.google_get(url, bypass_cloudflare=True)

    # Wait for page to fully load
    driver.long_random_sleep()
    # Locate iframe containing the Cloudflare challenge
    iframe = driver.get_element_at_point(160, 290)

    # Find checkbox element within the iframe
    checkbox = iframe.get_element_at_point(30, 30)

    # Enable human mode for realistic, human-like mouse movements
    driver.enable_human_mode()

    # Click the checkbox to solve the challenge
    checkbox.click()

    driver.disable_human_mode()

    # Now get the element using element selection methods
    content_tab = driver.wait_for_element("#chr-content")  # This finds the element

    if content_tab is None:
        print("Element #chr-content not found")
        return ""

    content_html = getattr(content_tab, "html", str(content_tab))

    # Extract only <p> tags
    soup = BeautifulSoup(content_html, "html.parser")
    chapter_text = "".join(str(p) for p in soup.find_all("p"))
    return chapter_text


@request(output=None, max_retry=10)
def scrape_chapter_with_request(request: Request, url: str) -> str:
    response = request.get(url, timeout=30)
    soup = soupify(response)

    if soup.get_text() == "Just a moment...Enable JavaScript and cookies to continue":
        return scrape_chapter(url)  # type: ignore

    # Content (keep <p> tags)
    el = soup.find(id="chr-content")
    chapter_text = (
        "".join(str(p) for p in el.find_all("p")) if isinstance(el, Tag) else ""
    )

    return chapter_text
