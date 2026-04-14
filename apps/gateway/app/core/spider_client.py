"""
Spider_XHS Data_Spider wrapper for HTTP Gateway
"""
import logging
import sys
from typing import Any, Tuple

from app.core import config, cookie_vault

logger = logging.getLogger(__name__)

# Discover Spider_XHS path
try:
    SPIDER_XHS_PATH = config.get_spider_xhs_path()
    sys.path.insert(0, str(SPIDER_XHS_PATH))
except config.StartupError:
    SPIDER_XHS_PATH = None


class DataSpiderClient:
    """Wrapper around Spider_XHS Data_Spider"""

    def __init__(self):
        self.cookies = cookie_vault.load_cookie()

    def spider_note(self, note_url: str) -> Tuple[bool, Any, dict]:
        """Collect single note with full info"""
        try:
            from apis.xhs_pc_apis import XHS_Apis
            from xhs_utils.data_util import handle_note_info

            xhs_apis = XHS_Apis()
            success, msg, note_info = xhs_apis.get_note_info(note_url, self.cookies)
            if success:
                note_info = note_info["data"]["items"][0]
                note_info["url"] = note_url
                processed = handle_note_info(note_info)
                return True, processed, processed
            return False, str(msg), {}
        except Exception as e:
            logger.exception("spider_note failed")
            return False, str(e), {}

    def spider_note_url_only(self, note_url: str) -> Tuple[bool, Any, dict]:
        """Collect single note info (raw, for batch processing)"""
        return self.spider_note(note_url)

    def spider_user_all_notes(self, user_url: str) -> Tuple[bool, Any, list]:
        """Collect all note URLs from a user"""
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
        """Search and collect note URLs"""
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
