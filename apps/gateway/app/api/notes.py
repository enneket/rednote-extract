"""
Notes collection endpoints
"""
import logging
from pathlib import Path

from fastapi import APIRouter, HTTPException
from pydantic import BaseModel

logger = logging.getLogger(__name__)

router = APIRouter()

SPIDER_XHS_PATH = Path(__file__).parent.parent.parent.parent.parent / "Spider_XHS"


class SingleNoteRequest(BaseModel):
    note_url: str


class UserNotesRequest(BaseModel):
    user_url: str


class SearchNotesRequest(BaseModel):
    keyword: str
    count: int = 10
    sort_type: int = 0  # 0:综合 1:最新 2:最多点赞 3:最多评论 4:最多收藏
    note_type: int = 0  # 0:不限 1:视频 2:普通
    note_time: int = 0  # 0:不限 1:一天内 2:一周内 3:半年内


@router.post("/single")
async def collect_single_note(req: SingleNoteRequest):
    if not req.note_url or "xiaohongshu.com" not in req.note_url:
        raise HTTPException(status_code=400, detail={"error": "invalid_url", "message": "invalid note_url format"})

    logger.info(f"Collecting single note: {req.note_url}")
    from app.core.spider_client import DataSpiderClient

    client = DataSpiderClient()
    success, msg, note_info = client.spider_note(req.note_url)

    if not success:
        raise HTTPException(status_code=500, detail={"error": "spider_error", "message": str(msg)})

    return note_info


@router.post("/user")
async def collect_user_notes(req: UserNotesRequest):
    if not req.user_url or "xiaohongshu.com" not in req.user_url:
        raise HTTPException(status_code=400, detail={"error": "invalid_url", "message": "invalid user_url format"})

    logger.info(f"Collecting user notes: {req.user_url}")
    from app.core.spider_client import DataSpiderClient

    client = DataSpiderClient()
    success, msg, note_list = client.spider_user_all_notes(req.user_url)

    if not success:
        raise HTTPException(status_code=500, detail={"error": "spider_error", "message": str(msg)})

    return {"notes": note_list, "count": len(note_list)}


@router.post("/search")
async def search_notes(req: SearchNotesRequest):
    if not req.keyword:
        raise HTTPException(status_code=400, detail="keyword is required")

    logger.info(f"Searching notes: {req.keyword}, count: {req.count}")
    from app.core.spider_client import DataSpiderClient

    client = DataSpiderClient()
    success, msg, note_list = client.spider_search(
        req.keyword, req.count, req.sort_type, req.note_type, req.note_time
    )

    if not success:
        raise HTTPException(status_code=500, detail={"error": "spider_error", "message": str(msg)})

    return {"notes": note_list, "count": len(note_list)}
