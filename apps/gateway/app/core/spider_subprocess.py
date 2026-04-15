"""
Spider_XHS subprocess wrapper.
Gateway spawns a minimal Python script that uses Spider_XHS as a pip package,
communicates via JSON stdin/stdout to avoid import conflicts.
"""
import json
import logging
import subprocess
import sys
import tempfile
from pathlib import Path
from typing import Any, Tuple

logger = logging.getLogger(__name__)

# Minimal script that runs Spider_XHS as pip package
# This avoids import conflicts because it runs in a fresh Python interpreter
SUBPROCESS_SCRIPT = '''
import sys
import json
import os

# Ensure Spider_XHS pip package is on path (or system Python)
try:
    from apis.xhs_pc_apis import XHS_Apis
    from xhs_utils.data_util import handle_note_info
except ImportError as e:
    print(json.dumps({"error": f"import_error: {e}"}))
    sys.exit(1)

def spider_note(note_url, cookies):
    xhs_apis = XHS_Apis()
    success, msg, info = xhs_apis.get_note_info(note_url, cookies)
    if success:
        info = info["data"]["items"][0]
        info["url"] = note_url
        info = handle_note_info(info)
        print(json.dumps({"success": True, "data": info}))
    else:
        print(json.dumps({"success": False, "error": str(msg)}))

def spider_user_all_notes(user_url, cookies):
    xhs_apis = XHS_Apis()
    success, msg, notes = xhs_apis.get_user_all_notes(user_url, cookies)
    if success:
        urls = []
        for n in notes:
            urls.append(f"https://www.xiaohongshu.com/explore/{n['note_id']}?xsec_token={n['xsec_token']}")
        print(json.dumps({"success": True, "data": urls}))
    else:
        print(json.dumps({"success": False, "error": str(msg)}))

def spider_search(keyword, count, sort_type, note_type, note_time, cookies):
    xhs_apis = XHS_Apis()
    success, msg, notes = xhs_apis.search_some_note(keyword, count, cookies, sort_type, note_type, note_time, 0, 0, None)
    if success:
        urls = []
        for n in notes:
            if n.get("model_type") == "note":
                urls.append(f"https://www.xiaohongshu.com/explore/{n['id']}?xsec_token={n['xsec_token']}")
        print(json.dumps({"success": True, "data": urls}))
    else:
        print(json.dumps({"success": False, "error": str(msg)}))

COMMANDS = {
    "spider_note": lambda p: spider_note(p["note_url"], p["cookies"]),
    "spider_user_all_notes": lambda p: spider_user_all_notes(p["user_url"], p["cookies"]),
    "spider_search": lambda p: spider_search(p["keyword"], p["count"], p.get("sort_type",0), p.get("note_type",0), p.get("note_time",0), p["cookies"]),
}

payload = json.loads(sys.argv[1])
cmd = payload.pop("__cmd__")
COMMANDS[cmd](payload)
'''


class SpiderSubprocess:
    """Wrapper that spawns Spider_XHS in a clean subprocess."""

    def __init__(self, cookies: str = ""):
        self.cookies = cookies

    def _run(self, cmd: str, payload: dict) -> dict:
        full_payload = {"__cmd__": cmd, "cookies": self.cookies, **payload}

        with tempfile.NamedTemporaryFile(
            mode="w", suffix=".json", delete=False
        ) as f:
            json.dump(full_payload, f)
            payload_file = f.name

        try:
            result = subprocess.run(
                [sys.executable, "-c", SUBPROCESS_SCRIPT, payload_file],
                capture_output=True,
                text=True,
                timeout=120,
            )
            if result.returncode != 0:
                logger.error(f"Subprocess stderr: {result.stderr[:300]}")
                return {"success": False, "error": result.stderr[:200]}

            output = result.stdout.strip()
            if not output:
                return {"success": False, "error": "empty output"}

            parsed = json.loads(output)
            if "error" in parsed:
                return {"success": False, "error": parsed["error"]}
            return {"success": parsed.get("success", False), "data": parsed.get("data", {})}

        except subprocess.TimeoutExpired:
            return {"success": False, "error": "timeout"}
        except Exception as e:
            logger.exception("Subprocess call failed")
            return {"success": False, "error": str(e)}
        finally:
            Path(payload_file).unlink(missing_ok=True)

    def spider_note(self, note_url: str) -> Tuple[bool, Any, dict]:
        result = self._run("spider_note", {"note_url": note_url})
        return result["success"], result.get("error", ""), result.get("data", {})

    def spider_user_all_notes(self, user_url: str) -> Tuple[bool, Any, list]:
        result = self._run("spider_user_all_notes", {"user_url": user_url})
        data = result.get("data", [])
        return result["success"], result.get("error", ""), data

    def spider_search(self, keyword: str, count: int, sort_type=0, note_type=0, note_time=0) -> Tuple[bool, Any, list]:
        result = self._run("spider_search", {
            "keyword": keyword, "count": count,
            "sort_type": sort_type, "note_type": note_type, "note_time": note_time,
        })
        data = result.get("data", [])
        return result["success"], result.get("error", ""), data
