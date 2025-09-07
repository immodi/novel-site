import re
from botasaurus.request import request, Request
from lxml import html
from scrapper import config


def slugify(text: str) -> str:
    """Make safe filenames"""
    return re.sub(r"[^a-zA-Z0-9_-]+", "_", text).strip("_")


def safe_text(el: html.HtmlElement) -> str | None:
    """Return stripped text or None"""
    if not hasattr(el, "text_content"):
        return None
    return el.text_content().strip()


def bold_green(text):
    return f"\033[1;32m{text}\033[0m"


@request(output=None)
def fetch_page(request: Request, url: str) -> html.HtmlElement:
    resp = request.get(url)
    resp.raise_for_status()
    return html.fromstring(resp.text)
