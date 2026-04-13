"""
Cookie management endpoints
"""
import json
import logging
from pathlib import Path

from fastapi import APIRouter, HTTPException
from pydantic import BaseModel

logger = logging.getLogger(__name__)

router = APIRouter()

COOKIE_FILE = Path(__file__).parent.parent.parent.parent / "Spider_XHS" / ".cookies"


class CookiesSetRequest(BaseModel):
    cookies: str


class CookiesGetResponse(BaseModel):
    configured: bool


@router.get("", response_model=CookiesGetResponse)
async def get_cookies():
    if not COOKIE_FILE.exists():
        return CookiesGetResponse(configured=False)
    try:
        COOKIE_FILE.read_text()
        return CookiesGetResponse(configured=True)
    except Exception:
        return CookiesGetResponse(configured=False)


@router.post("")
async def set_cookies(req: CookiesSetRequest):
    if not req.cookies:
        raise HTTPException(status_code=400, detail="cookies cannot be empty")

    try:
        COOKIE_FILE.parent.mkdir(parents=True, exist_ok=True)
        COOKIE_FILE.write_text(req.cookies)
        logger.info("Cookies updated successfully")
        return {"message": "Cookies saved"}
    except Exception as e:
        logger.error(f"Failed to save cookies: {e}")
        raise HTTPException(status_code=500, detail=f"Failed to save cookies: {e}")
