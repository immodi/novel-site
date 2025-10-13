from bs4 import Tag, BeautifulSoup
from botasaurus.request import request, Request
from botasaurus.soupify import soupify
from botasaurus.browser import browser, Driver
from urllib.parse import urljoin
from scrapper.helpers import utils, saver, helpers
from scrapper.modules.factories.factory import Parser, SkipDuplicate
from scrapper.helpers.utils import safe_text
from scrapper import config
from typing import List, Union, Generator
from scrapper.datatypes.novel import ChapterData, NovelData, NovelLink
from lxml import html
from scrapper.datatypes.novel_fire_chapter import NovelFireChapter
from typing import Tuple
from scrapper.cache.db_cache import NovelDataCache


class NovelFireParser(Parser):
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
            for novel_el in t.cssselect(".novel-list > *"):
                title_el = novel_el.cssselect("h4")
                link_el = novel_el.cssselect("a")

                if not title_el or not link_el:
                    continue

                title = title_el[0].text_content().strip()
                url = urljoin(config.BASE_URL, link_el[0].get("href"))

                if (
                    self.skip_duplicates == SkipDuplicate.NOVEL
                    and self.cache.novel_exists(url)
                ):
                    print(f"Novel {title} already exists. Skipping...")
                    continue

                novels.append(NovelLink(title=title, url=url))

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
        genres = (
            [g.text_content().strip() for g in genres_el] if len(genres_el) > 0 else []
        )

        # Tags
        tags_el = tree.cssselect("ul.content li")
        tags = [t.text_content().strip() for t in tags_el] if len(tags_el) > 0 else []

        # Cover image
        cover_el = tree.cssselect(".glass-background img")
        cover_image = cover_el[0].get("src") if len(cover_el) > 0 else ""

        # Description
        desc_el = tree.cssselect(".content p")
        if len(desc_el) > 0:
            # Get the text of all <p> tags inside .content, join with newline
            description = "\n".join([p.text_content().strip() for p in desc_el])
        else:
            description = ""

        novel_data = NovelData(
            title=title or "",
            author=author,
            genres=genres,
            status=status,
            tags=tags,
            cover_image=cover_image,
            description=description,
            url=url,
        )

        saver.save_item(novel_data, f"{config.OUTPUT_DIR}/novels")
        if save_image and cover_image:
            helpers.download_image(
                cover_image, f"{config.OUTPUT_DIR}/covers", title or "cover"
            )

        # self.cache.save_novel(novel_data)

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
            # self.cache.remove_novel_by_url(url)
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

            chapter = ChapterData(
                title=chapter_data.title,
                content=chapter_data.content,
                url=chapter_link,
            )

            chapters.append(chapter)

            if save_per_chapter:
                saver.save_item(
                    chapter, f"{config.OUTPUT_DIR}/chapters/{utils.slugify(novel_name)}"
                )
                # self.cache.save_last_chapter(url, chapter.url, chapter.title)

            print(
                f"--> Fetched {chapter_data.title} of length {utils.bold_green(len(chapter_data.content))} from url {chapter_link}."
            )

        if not save_per_chapter:
            for chapter in chapters:
                saver.save_item(
                    chapter, f"{config.OUTPUT_DIR}/chapters/{utils.slugify(novel_name)}"
                )
            # self.cache.save_last_chapter(url, chapters[-1].url, chapters[-1].title)

        return chapters

    def parse_chapters_with_notify(
        self, url: str, novel_name: str, save_per_chapter: bool
    ) -> Generator[str, None, List[ChapterData]]:
        """
        Like parse_chapters, but yields a message for each chapter fetched.
        Returns the full list of chapters at the end.
        """
        ajax_url = f"{url}/chapters"
        tree = utils.fetch_page(ajax_url)  # type: ignore
        chapters: List[ChapterData] = []

        # Find last chapter number
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
            yield f"Maximum number of chapters ({self.max_chapters_number}) exceeded. Skipping {novel_name}."
            self.clean_up_novel(novel_name)
            return chapters

        # Fetch chapters
        for i in range(1, last_chapter_number + 1):
            chapter_link = f"{base_url}chapter-{i}"
            chapter_data = scrape_chapter_with_request(chapter_link)  # type: ignore

            if self.skip_duplicates == SkipDuplicate.CHAPTER and self.chapter_exists(
                chapter_data.title, utils.slugify(novel_name)
            ):
                yield f"Chapter {chapter_data.title} already exists. Skipping..."
                continue

            chapter = ChapterData(
                title=chapter_data.title,
                content=chapter_data.content,
                url=chapter_link,
            )
            chapters.append(chapter)

            if save_per_chapter:
                saver.save_item(
                    chapter, f"{config.OUTPUT_DIR}/chapters/{utils.slugify(novel_name)}"
                )

            msg = f"--> Fetched {chapter.title} of length {utils.bold_green(len(chapter.content))} from url {chapter_link}."
            yield msg

        if not save_per_chapter:
            for chapter in chapters:
                saver.save_item(
                    chapter, f"{config.OUTPUT_DIR}/chapters/{utils.slugify(novel_name)}"
                )

        return chapters

    def update_novel_with_notify(
        self, novel_name: str, novel_url: str, last_chapter_url: str
    ) -> Generator[str, None, Tuple[str, List[ChapterData]]]:
        chapters: List[ChapterData] = []
        if not last_chapter_url:
            yield f"No last chapter found for novel {novel_name}"

        current_url = last_chapter_url
        yield f"Updating {novel_name} from {current_url}"
        save_dir = f"{config.OUTPUT_DIR}/chapters/{utils.slugify(novel_name)}/updates"

        while True:
            # 1. request the page
            tree: html.HtmlElement = utils.fetch_page(current_url)  # type: ignore

            # 2. get the "nextchap" element
            next_el = tree.cssselect(".nextchap")
            if not next_el:
                yield f"No .nextchap element found at {current_url}"
                break

            next_btn = next_el[0]

            # 3. if it has 'isDisabled' class, stop
            classes = next_btn.get("class", "")
            if "isDisabled" in classes.split():
                print("Reached the final chapter.")
                break

            # 4. follow the next link
            next_href = next_btn.get("href")
            if not next_href:
                yield f"No href found for next chapter at {current_url}"
                break

            # Make sure the href is absolute if needed
            if next_href.startswith("/"):
                # adjust base if needed; assuming same site as current_url
                base = current_url.split("/chapter-")[0]
                next_href = base.rstrip("/") + next_href

            current_url = next_href

            # 5. scrape chapter content
            chapter_data = scrape_chapter_with_request(current_url)  # type: ignore

            chapter = ChapterData(
                title=chapter_data.title,
                content=chapter_data.content,
                url=current_url,
            )
            chapters.append(chapter)

            # save each chapter immediately
            saver.save_item(chapter, save_dir)
            self.cache.save_last_chapter(novel_url, current_url, chapter.title)
            msg = f"""--> Fetched {chapter_data.title} of length \n 
            {utils.bold_green(len(chapter_data.content))} from url {current_url}."""
            yield msg

        return save_dir, chapters

    def update_novel(
        self, novel_name: str, novel_url: str, last_chapter_url: str
    ) -> Tuple[str, List[ChapterData]]:
        chapters: List[ChapterData] = []
        if not last_chapter_url:
            print("No last chapter found.")
            return "", chapters

        current_url = last_chapter_url
        print(f"Updating {novel_name} from {current_url}")
        save_dir = f"{config.OUTPUT_DIR}/chapters/{utils.slugify(novel_name)}/updates"

        while True:
            # 1. request the page
            tree: html.HtmlElement = utils.fetch_page(current_url)  # type: ignore

            # 2. get the "nextchap" element
            next_el = tree.cssselect(".nextchap")
            if not next_el:
                print(f"No .nextchap element found at {current_url}")
                break

            next_btn = next_el[0]

            # 3. if it has 'isDisabled' class, stop
            classes = next_btn.get("class", "")
            if "isDisabled" in classes.split():
                print("Reached the final chapter.")
                break

            # 4. follow the next link
            next_href = next_btn.get("href")
            if not next_href:
                print(f"No href found for next chapter at {current_url}")
                break

            # Make sure the href is absolute if needed
            if next_href.startswith("/"):
                # adjust base if needed; assuming same site as current_url
                base = current_url.split("/chapter-")[0]
                next_href = base.rstrip("/") + next_href

            current_url = next_href

            # 5. scrape chapter content
            chapter_data = scrape_chapter_with_request(current_url)  # type: ignore

            chapter = ChapterData(
                title=chapter_data.title,
                content=chapter_data.content,
                url=current_url,
            )
            chapters.append(chapter)

            # save each chapter immediately
            saver.save_item(chapter, save_dir)
            self.cache.save_last_chapter(novel_url, current_url, chapter.title)

            print(
                f"--> Fetched {chapter_data.title} of length "
                f"{utils.bold_green(len(chapter_data.content))} from url {current_url}."
            )

        return save_dir, chapters


@browser(output=None, headless=True, max_retry=10)
def scrape_chapter(driver: Driver, url: str) -> NovelFireChapter:
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

    # Title
    title = driver.get_text(".chapter-title").strip()

    # Get raw HTML of content container
    content_tab = driver.get("#content")
    content_html = getattr(content_tab, "html", str(content_tab))  # fallback to str

    # Extract only <p> tags
    soup = BeautifulSoup(content_html, "html.parser")

    content = "".join(str(p) for p in soup.find_all("p"))

    return NovelFireChapter(title=title, content=content)


@request(output=None, max_retry=10)
def scrape_chapter_with_request(request: Request, url: str) -> NovelFireChapter:
    response = request.get(url, timeout=30)
    soup = soupify(response)

    if soup.get_text() == "Just a moment...Enable JavaScript and cookies to continue":
        return scrape_chapter(url)  # type: ignore

    # Title
    chapter_title = soup.find("span", class_="chapter-title")
    title = chapter_title.get_text(strip=True) if chapter_title else ""
    title = (title[:197] + "...") if len(title) > 200 else title

    # Content (keep <p> tags)
    chapter_content = soup.find(id="content")
    content = (
        "".join(str(p) for p in chapter_content.find_all("p"))
        if isinstance(chapter_content, Tag)
        else ""
    )

    return NovelFireChapter(title=title, content=content)
