"""
Spider_XHS Data_Spider wrapper for HTTP Gateway
Directly imports Spider_XHS modules instead of subprocess.
"""
import json
import logging
import sys
from pathlib import Path
from typing import Any, Tuple

logger = logging.getLogger(__name__)

# Add Spider_XHS to path
SPIDER_XHS_PATH = Path(__file__).parent.parent.parent.parent.parent / "Spider_XHS"
sys.path.insert(0, str(SPIDER_XHS_PATH))

COOKIE_FILE = SPIDER_XHS_PATH / ".cookies"


class DataSpiderClient:
    """Wrapper around Spider_XHS Data_Spider"""

    def __init__(self):
        self.cookies = self._load_cookies()

    def _load_cookies(self) -> str:
        if not COOKIE_FILE.exists():
            return ""
        return COOKIE_FILE.read_text().strip()

    def spider_note(self, note_url: str) -> Tuple[bool, Any, dict]:
        """Collect single note"""
        try:
            from apis.xhs_pc_apis import XHS_Apis
            from xhs_utils.data_util import handle_note_info

            xhs_apis = XHS_Apis()
            success, msg, note_info = xhs_apis.get_note_info(note_url, self.cookies)
            if success:
                note_info = note_info["data"]["items"][0]
                note_info["url"] = note_url
                note_info = handle_note_info(note_info)
                return True, note_info, note_info
            return False, str(msg), {}
        except Exception as e:
            logger.exception("spider_note failed")
            return False, str(e), {}

    def spider_user_all_notes(self, user_url: str) -> Tuple[bool, Any, list]:
        """Collect all notes from a user"""
        try:
            from apis.xhs_pc_apis import XHS_Apis

            xhs_apis = XHS_Apis()
            success, msg, all_notes = xhs_apis.get_user_all_notes(user_url, self.cookies)
            if success:
                note_urls = []
                for note in all_notes:
                    note_url = (
                        f"https://www.xiaohongshu.com/explore/{note['note_id']}"
                        f"?xsec_token={note['xsec_token']}"
                    )
                    note_urls.append(note_url)
                return True, note_urls, note_urls
            return False, str(msg), []
        except Exception as e:
            logger.exception("spider_user_all_notes failed")
            return False, str(e), []

    def spider_search(
        self, keyword: str, count: int, sort_type: int = 0,
        note_type: int = 0, note_time: int = 0
    ) -> Tuple[bool, Any, list]:
        """Search and collect notes"""
        try:
            from apis.xhs_pc_apis import XHS_Apis

            xhs_apis = XHS_Apis()
            success, msg, notes = xhs_apis.search_some_note(
                keyword, count, self.cookies, sort_type, note_type, note_time, 0, 0, None
            )
            if success:
                note_urls = []
                for note in notes:
                    if note.get("model_type") == "note":
                        note_url = (
                            f"https://www.xiaohongshu.com/explore/{note['id']}"
                            f"?xsec_token={note['xsec_token']}"
                        )
                        note_urls.append(note_url)
                return True, note_urls, note_urls
            return False, str(msg), []
        except Exception as e:
            logger.exception("spider_search failed")
            return False, str(e), []
