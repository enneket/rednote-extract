"""
Spider_XHS wrapper - uses subprocess to avoid import conflicts.
Spider_XHS must be installed as pip package: pip install spider-xhs
Or set SPIDER_XHS_PATH to a pip-installable Spider_XHS.
"""
from typing import Any, Tuple

from app.core import cookie_vault
from app.core.spider_subprocess import SpiderSubprocess


class DataSpiderClient:
    """Wrapper around Spider_XHS Data_Spider via subprocess."""

    def __init__(self):
        self.cookies = cookie_vault.load_cookie()

    def spider_note(self, note_url: str) -> Tuple[bool, Any, dict]:
        """Collect single note with full info."""
        client = SpiderSubprocess(self.cookies)
        return client.spider_note(note_url)

    def spider_user_all_notes(self, user_url: str) -> Tuple[bool, Any, list]:
        """Collect all note URLs from a user."""
        client = SpiderSubprocess(self.cookies)
        return client.spider_user_all_notes(user_url)

    def spider_search(
        self, keyword: str, count: int, sort_type: int = 0,
        note_type: int = 0, note_time: int = 0
    ) -> Tuple[bool, Any, list]:
        """Search and collect note URLs."""
        client = SpiderSubprocess(self.cookies)
        return client.spider_search(keyword, count, sort_type, note_type, note_time)
