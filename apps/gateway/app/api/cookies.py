"""
Cookie management endpoints
"""
import logging

from fastapi import APIRouter, HTTPException
from pydantic import BaseModel

from app.core import cookie_vault

logger = logging.getLogger(__name__)

router = APIRouter()


class CookiesSetRequest(BaseModel):
    cookies: str


class CookiesGetResponse(BaseModel):
    configured: bool
    has_content: bool
    domain: str = "xiaohongshu.com"


@router.get("", response_model=CookiesGetResponse)
async def get_cookies():
    return CookiesGetResponse(
        configured=cookie_vault.is_configured(),
        has_content=cookie_vault.has_content(),
        domain="xiaohongshu.com",
    )


@router.post("")
async def set_cookies(req: CookiesSetRequest):
    if not req.cookies:
        raise HTTPException(status_code=400, detail="cookies cannot be empty")

    try:
        cookie_vault.save_cookie(req.cookies)
        logger.info("Cookies saved to vault")
        return {"message": "Cookies saved"}
    except Exception as e:
        logger.error(f"Failed to save cookies: {e}")
        raise HTTPException(status_code=500, detail=f"Failed to save cookies: {e}")


@router.delete("")
async def delete_cookies():
    try:
        cookie_vault.delete_cookie()
        return {"message": "Cookies deleted"}
    except Exception as e:
        logger.error(f"Failed to delete cookies: {e}")
        raise HTTPException(status_code=500, detail=f"Failed to delete cookies: {e}")
