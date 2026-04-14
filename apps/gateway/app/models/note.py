"""
Unified note data model for API responses
"""
from typing import List, Optional

from pydantic import BaseModel


class NoteData(BaseModel):
    """Unified note data structure returned by all collection APIs."""
    id: str = ""
    title: str = ""
    author: str = ""
    author_id: str = ""
    type: str = ""  # 图集 / 视频
    liked: int = 0
    collected: int = 0
    commented: int = 0
    shared: int = 0
    url: str = ""
    tags: List[str] = []
    time: str = ""  # e.g. "2024-01-01 12:00:00"
    ip_location: str = ""


def map_handle_note_info_to_note_data(raw: dict) -> NoteData:
    """Map Spider_XHS handle_note_info output to NoteData."""
    return NoteData(
        id=raw.get("id", ""),
        title=raw.get("title", ""),
        author=raw.get("nickname", ""),
        author_id=raw.get("user_id", ""),
        type=raw.get("type", ""),
        liked=_safe_int(raw.get("liked_count", 0)),
        collected=_safe_int(raw.get("collected_count", 0)),
        commented=_safe_int(raw.get("comment_count", 0)),
        shared=_safe_int(raw.get("share_count", 0)),
        url=raw.get("url", ""),
        tags=raw.get("tags", []),
        time=raw.get("time", ""),
        ip_location=raw.get("ip_location", ""),
    )


def _safe_int(val) -> int:
    if isinstance(val, int):
        return val
    if isinstance(val, str) and val.isdigit():
        return int(val)
    try:
        return int(val)
    except (ValueError, TypeError):
        return 0
