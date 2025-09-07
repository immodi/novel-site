from scrapper.modules.factories.interface import Parser, SkipDuplicate
from scrapper.modules.parsers.novelbin_parser import NovelBinParser
from scrapper.modules.parsers.novelfire_parser import NovelFireParser


def get_parser(
    url: str, skip_duplicates: SkipDuplicate, max_chapters_number: int
) -> Parser:
    """
    Returns the appropriate parser based on config.LIST_TREE content.
    Defaults to NovelBinParser if no match is found.
    """
    if "novelbin" in url.lower():
        return NovelBinParser(max_chapters_number, skip_duplicates)
    elif "novelfire" in url.lower():
        return NovelFireParser(max_chapters_number, skip_duplicates)

    # default fallback
    return NovelBinParser(max_chapters_number)
