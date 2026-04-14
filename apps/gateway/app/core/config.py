"""
Gateway configuration and Spider_XHS path discovery
"""
import importlib.util
import logging
import os
import sys
from pathlib import Path

logger = logging.getLogger(__name__)

# Gateway config directory
GATEWAY_HOME = Path.home() / ".xhs-gateway"
CONFIG_FILE = GATEWAY_HOME / "config.yaml"


def _get_bundle_dir() -> Path:
    """Get the directory where the frozen exe is extracted (PyInstaller)."""
    if getattr(sys, "frozen", False) and hasattr(sys, "_MEIPASS"):
        return Path(sys._MEIPASS)
    if getattr(sys, "frozen", False):
        return Path(sys.executable).parent
    return None


def get_spider_xhs_path() -> Path:
    """
    Discover Spider_XHS path using priority:
    1. SPIDER_XHS_PATH environment variable
    2. Bundled vendor/Spider_XHS (PyInstaller _MEIPASS)
    3. pip-installed package
    4. sibling ../Spider_XHS directory (development)
    5. raise StartupError
    """
    # 1. Environment variable
    env_path = os.environ.get("SPIDER_XHS_PATH")
    if env_path:
        path = Path(env_path).resolve()
        if path.exists():
            logger.info(f"Using Spider_XHS from SPIDER_XHS env: {path}")
            return path
        logger.warning(f"SPIDER_XHS_PATH set but not found: {path}")

    # 2. Bundled vendor path (PyInstaller)
    bundle_dir = _get_bundle_dir()
    if bundle_dir:
        vendor_spider = bundle_dir / "vendor" / "Spider_XHS"
        if vendor_spider.exists():
            logger.info(f"Using Spider_XHS from bundled vendor: {vendor_spider}")
            return vendor_spider

    # 3. Pip package
    spec = importlib.util.find_spec("spider_xhs")
    if spec and spec.submodule_search_locations:
        path = Path(spec.submodule_search_locations[0]).parent
        logger.info(f"Using Spider_XHS from pip package: {path}")
        return path

    # 4. Sibling directory (development)
    sibling = Path(__file__).parent.parent.parent.parent / "Spider_XHS"
    if sibling.exists():
        logger.info(f"Using Spider_XHS from sibling directory: {sibling}")
        return sibling

    raise StartupError(
        "Spider_XHS not found. Set SPIDER_XHS_PATH environment variable, "
        "install via 'pip install spider-xhs', or ensure vendor/Spider_XHS is bundled."
    )


def ensure_gateway_home() -> Path:
    """Ensure ~/.xhs-gateway/ directory and config exist."""
    GATEWAY_HOME.mkdir(parents=True, exist_ok=True)
    return GATEWAY_HOME


class StartupError(Exception):
    """Raised when Gateway cannot start due to missing dependencies."""
    pass
