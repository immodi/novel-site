import json
from dataclasses import fields, is_dataclass
from typing import Type, TypeVar

T = TypeVar("T")


def load_from_json(path: str, cls: Type[T]) -> T:
    """Load any dataclass object from a JSON file."""
    if not is_dataclass(cls):
        raise ValueError(f"{cls} must be a dataclass type")

    with open(path, "r", encoding="utf-8") as f:
        data = json.load(f)

    # Only keep fields that exist in the dataclass
    cls_fields = {f.name for f in fields(cls)}
    filtered_data = {k: v for k, v in data.items() if k in cls_fields}

    return cls(**filtered_data)
