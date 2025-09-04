from urllib.parse import urljoin
from botasaurus.task import task
from botasaurus.request import request, Request
from botasaurus.soupify import soupify
from botasaurus.browser import browser, Driver
from modules import utils
from .interface import Parser
from modules.utils import safe_text
from modules import config
from typing import List
from datatypes.novel import ChapterData, NovelData, NovelLink
from lxml import html


class NovelBinParser(Parser):
    def parse_list_of_novels(self, tree: html.HtmlElement) -> List[NovelLink]:
        """Extract novel links from the list page"""
        return [
            {
                "title": a.text_content().strip(),
                "url": urljoin(config.BASE_URL, a.get("href")),
            }
            for a in tree.cssselect(".list-novel .row .novel-title > a")
        ]

    def parse_novel(self, tree: html.HtmlElement, url: str) -> NovelData:
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
        description_tree = utils.fetch_page(description_url)
        desc_el = description_tree.cssselect(".desc-text")
        description = safe_text(desc_el[0]) if desc_el else None

        return {
            "title": title,
            "author": author,
            "genres": genres,
            "status": status,
            "tags": tags,
            "cover_image": img,
            "description": description,
            "url": url,
        }

    def parse_chapters(self, url: str) -> List[ChapterData]:
        """Fetch chapters using the first link and navigating via #next_chap button"""
        # Fetch chapter list via AJAX
        novel_id = url.rstrip("/").split("/")[-1]
        ajax_url = f"{config.BASE_URL}/ajax/chapter-archive?novelId={novel_id}"
        tree: html.HtmlElement = utils.fetch_page(ajax_url)

        chapters: List[ChapterData] = []

        chapters_links = tree.cssselect(".panel-body a")
        for chapter_link in chapters_links:
            chapter_title = chapter_link.text_content().strip()
            chapter_url = urljoin(config.BASE_URL, chapter_link.get("href"))
            chapter_content = scrape_chapter_with_request(chapter_url)  # type: ignore

            chapters.append(
                {"title": chapter_title, "content": chapter_content, "url": chapter_url}
            )

            print(f"--> Fetched {chapter_title} of url {chapter_url}.")

        return chapters


@browser
def scrape_chapter(driver: Driver, url: str) -> str:
    driver.get(url, True)
    chapter_text = driver.get_text("#chr-content")
    return chapter_text


@request(output=None)
def scrape_chapter_with_request(request: Request, url: str) -> str:
    response = request.get(url)
    soup = soupify(response)
    el = soup.find(id="chr-content")
    chapter_text = el.get_text().strip() if el else ""
    return chapter_text
