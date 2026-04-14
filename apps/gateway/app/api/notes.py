"""
Notes collection endpoints - unified response with full note data
"""
import logging
from typing import List, Optional

from fastapi import APIRouter, HTTPException, Query
from pydantic import BaseModel

from app.core.spider_client import DataSpiderClient
from app.models.note import NoteData, map_handle_note_info_to_note_data

logger = logging.getLogger(__name__)

router = APIRouter()


class SingleNoteRequest(BaseModel):
    note_url: str


class UserNotesRequest(BaseModel):
    user_url: str


class SearchNotesRequest(BaseModel):
    keyword: str
    count: int = 10
    sort_type: int = 0
    note_type: int = 0
    note_time: int = 0


class NotesListResponse(BaseModel):
    notes: List[NoteData]
    count: int
    partial: bool = False
    errors: List[str] = []


class SingleNoteResponse(BaseModel):
    note: NoteData
    partial: bool = False


@router.post("/single", response_model=SingleNoteResponse)
async def collect_single_note(req: SingleNoteRequest):
    if not req.note_url or "xiaohongshu.com" not in req.note_url:
        raise HTTPException(status_code=400, detail={"error": "invalid_url", "message": "invalid note_url format"})

    logger.info(f"Collecting single note: {req.note_url}")
    client = DataSpiderClient()
    success, msg, note_data = client.spider_note(req.note_url)

    if not success:
        raise HTTPException(status_code=500, detail={"error": "spider_error", "message": str(msg)})

    mapped = map_handle_note_info_to_note_data(note_data)
    return SingleNoteResponse(note=mapped)


@router.post("/user", response_model=NotesListResponse)
async def collect_user_notes(req: UserNotesRequest):
    if not req.user_url or "xiaohongshu.com" not in req.user_url:
        raise HTTPException(status_code=400, detail={"error": "invalid_url", "message": "invalid user_url format"})

    logger.info(f"Collecting user notes: {req.user_url}")
    client = DataSpiderClient()
    success, msg, note_urls = client.spider_user_all_notes(req.user_url)

    if not success:
        raise HTTPException(status_code=500, detail={"error": "spider_error", "message": str(msg)})

    return _fetch_notes_batch(note_urls)


@router.post("/search", response_model=NotesListResponse)
async def search_notes(req: SearchNotesRequest):
    if not req.keyword:
        raise HTTPException(status_code=400, detail="keyword is required")

    logger.info(f"Searching notes: {req.keyword}, count: {req.count}")
    client = DataSpiderClient()
    success, msg, note_urls = client.spider_search(
        req.keyword, req.count, req.sort_type, req.note_type, req.note_time
    )

    if not success:
        raise HTTPException(status_code=500, detail={"error": "spider_error", "message": str(msg)})

    return _fetch_notes_batch(note_urls)


def _fetch_notes_batch(note_urls: List[str]) -> NotesListResponse:
    """Fetch full info for each note URL, with graceful degradation."""
    notes = []
    errors = []

    for url in note_urls:
        try:
            client = DataSpiderClient()
            success, msg, raw = client.spider_note(url)
            if success:
                mapped = map_handle_note_info_to_note_data(raw)
                notes.append(mapped)
            else:
                errors.append(f"{url}: {msg}")
        except Exception as e:
            errors.append(f"{url}: {e}")

    return NotesListResponse(
        notes=notes,
        count=len(notes),
        partial=len(errors) > 0,
        errors=errors,
    )
