"""
XHS HTTP Gateway - FastAPI service wrapping Spider_XHS
"""
import logging
from contextlib import asynccontextmanager

from dotenv import load_dotenv
from fastapi import FastAPI
from fastapi.middleware.cors import CORSMiddleware

from app.api import cookies, health, notes, export
from app.core import config

# Load .env
load_dotenv()

logging.basicConfig(level=logging.INFO)
logger = logging.getLogger(__name__)


@asynccontextmanager
async def lifespan(app: FastAPI):
    logger.info("XHS Gateway starting up...")
    # Ensure gateway home dir
    config.ensure_gateway_home()
    # Discover Spider_XHS (fail fast)
    try:
        spider_path = config.get_spider_xhs_path()
        app.state.spider_xhs_path = spider_path
        logger.info(f"Spider_XHS ready at: {spider_path}")
    except config.StartupError as e:
        logger.error(f"Startup failed: {e}")
        raise
    yield
    logger.info("XHS Gateway shutting down...")


app = FastAPI(
    title="XHS HTTP Gateway",
    description="HTTP API wrapper for Spider_XHS爬虫",
    version="1.0.0",
    lifespan=lifespan,
)

# CORS for desktop client
app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)

# Register routes
app.include_router(health.router, tags=["health"])
app.include_router(cookies.router, prefix="/api/v1/cookies", tags=["cookies"])
app.include_router(notes.router, prefix="/api/v1/notes", tags=["notes"])
app.include_router(export.router, prefix="/api/v1/notes", tags=["notes"])


if __name__ == "__main__":
    import sys
    import uvicorn
    is_frozen = getattr(sys, "frozen", False)
    uvicorn.run("main:app", host="127.0.0.1", port=8000, reload=not is_frozen)
