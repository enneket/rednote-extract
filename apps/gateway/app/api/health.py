"""
Health check endpoint
"""
from fastapi import APIRouter
from pydantic import BaseModel

router = APIRouter()


class HealthResponse(BaseModel):
    status: str
    spider_ready: bool


@router.get("/health", response_model=HealthResponse)
async def health_check():
    return HealthResponse(status="ok", spider_ready=True)
