from pathlib import Path
import shutil


def move_files_up_and_delete(dir_path: str) -> None:
    """
    Move every file (and subfolder) inside `dir_path` to its direct parent
    directory, then remove `dir_path` itself.

    Parameters
    ----------
    dir_path : str
        Path of the directory whose contents should be moved up one level.

    Raises
    ------
    FileNotFoundError
        If `dir_path` does not exist or is not a directory.
    """
    src = Path(dir_path).resolve()

    if not src.exists() or not src.is_dir():
        raise FileNotFoundError(f"Directory not found: {src}")

    parent = src.parent

    for item in src.iterdir():
        # build target path inside the parent directory
        dest = parent / item.name

        # if a file/dir with the same name already exists, choose behaviour:
        if dest.exists():
            raise FileExistsError(f"Destination already exists: {dest}")

        # move the file or directory
        shutil.move(str(item), str(parent))

    # finally remove the now-empty directory itself
    src.rmdir()
