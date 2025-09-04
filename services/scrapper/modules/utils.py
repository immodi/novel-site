import re
import requests
from lxml import html
from . import config


def slugify(text: str) -> str:
    """Make safe filenames"""
    return re.sub(r"[^a-zA-Z0-9_-]+", "_", text).strip("_")


def safe_text(el: html.HtmlElement) -> str | None:
    """Return stripped text or None"""
    if not hasattr(el, "text_content"):
        return None
    return el.text_content().strip()


def fetch_page(url: str) -> html.HtmlElement:
    resp = requests.get(url, headers=config.HEADERS)
    resp.raise_for_status()
    return html.fromstring(resp.text)
