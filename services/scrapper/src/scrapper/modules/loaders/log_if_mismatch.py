import logging
from pathlib import Path
from scrapper.datatypes.get_last_chapter import GetLastChapterResponse


def log_if_mismatch(
    remote: GetLastChapterResponse,
    local_name: str,
    novel_identifier: str | int,
) -> bool:
    """
    Compare remote and local last-chapter info and:
      • create /logs folder if missing
      • log mismatch into /logs/novel_{novel_name}_chapter_mismatch.log
      • return True if mismatch found, else False
    """
    logs_dir = Path("logs")
    logs_dir.mkdir(parents=True, exist_ok=True)

    log_file = logs_dir / f"novel_{novel_identifier}_chapter_mismatch.log"

    # Configure a logger that writes to logs/novel_{}_chapter_mismatch.log
    logger = logging.getLogger(f"novel_{novel_identifier}_chapter_mismatch")
    if not logger.handlers:  # avoid adding multiple handlers on repeated calls
        logger.setLevel(logging.INFO)
        handler = logging.FileHandler(log_file, encoding="utf-8")
        handler.setFormatter(
            logging.Formatter("%(asctime)s [%(levelname)s] %(message)s")
        )
        logger.addHandler(handler)

    mismatch = remote.last_chapter_name != local_name

    if mismatch:
        logger.warning(
            "Mismatch for novel %s (remote ID %s): "
            "Local ->  name=%r | "
            "Remote -> name=%r",
            novel_identifier,
            remote.novel_id,
            local_name,
            remote.last_chapter_name,
        )

    return mismatch
