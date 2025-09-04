from modules.parsers.interface import Parser
from modules.parsers.novelbinparser import NovelBinParser


def get_parser(url: str) -> Parser:
    """
    Returns the appropriate parser based on config.LIST_TREE content.
    Defaults to NovelBinParser if no match is found.
    """
    if "novelbin" in url.lower():
        return NovelBinParser()
    # elif "other-site" in list_tree_content:
    #     return OtherParser()

    # default fallback
    return NovelBinParser()
