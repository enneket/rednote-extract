"""
Cookie vault: encrypted storage for Xiaohongshu cookies.
Uses AES-256-GCM with a machine-specific key derived from hostname + username.
"""
import base64
import getpass
import hashlib
import json
import logging
import socket
import shutil
from pathlib import Path

from cryptography.hazmat.backends import default_backend
from cryptography.hazmat.primitives.ciphers.aead import AESGCM

from app.core.config import GATEWAY_HOME

logger = logging.getLogger(__name__)

COOKIE_FILE = GATEWAY_HOME / "cookies.enc"
LEGACY_COOKIE = Path(__file__).parent.parent.parent.parent / "Spider_XHS" / ".cookies"


def _derive_key() -> bytes:
    """Derive a machine-specific key from hostname + username."""
    info = f"{socket.gethostname()}:{getpass.getuser()}:xhs-gateway-v1"
    return hashlib.sha256(info.encode()).digest()


def _encrypt(plaintext: str) -> bytes:
    """Encrypt plaintext with AES-256-GCM."""
    key = _derive_key()
    aesgcm = AESGCM(key)
    nonce = __import__("os").urandom(12)
    ciphertext = aesgcm.encrypt(nonce, plaintext.encode(), None)
    return base64.b64encode(nonce + ciphertext)


def _decrypt(data: bytes) -> str:
    """Decrypt ciphertext with AES-256-GCM."""
    key = _derive_key()
    aesgcm = AESGCM(key)
    raw = base64.b64decode(data)
    nonce, ciphertext = raw[:12], raw[12:]
    return aesgcm.decrypt(nonce, ciphertext, None).decode()


def migrate_legacy_cookie() -> bool:
    """Migrate cookie from Spider_XHS/.cookies to vault if needed."""
    if not LEGACY_COOKIE.exists():
        return False
    if COOKIE_FILE.exists():
        # Already migrated, remove legacy
        LEGACY_COOKIE.unlink()
        logger.info("Legacy cookie removed (vault already exists)")
        return True

    try:
        content = LEGACY_COOKIE.read_text().strip()
        if content:
            save_cookie(content)
            LEGACY_COOKIE.unlink()
            logger.info("Cookies migrated from Spider_XHS/.cookies to vault")
            return True
    except Exception as e:
        logger.warning(f"Cookie migration failed: {e}")
    return False


def save_cookie(cookie_str: str) -> None:
    """Encrypt and save cookie string."""
    COOKIE_FILE.parent.mkdir(parents=True, exist_ok=True)
    encrypted = _encrypt(cookie_str)
    COOKIE_FILE.write_bytes(encrypted)
    logger.info("Cookies saved to vault")


def load_cookie() -> str:
    """Load and decrypt cookie string."""
    if not COOKIE_FILE.exists():
        return ""
    try:
        return _decrypt(COOKIE_FILE.read_bytes())
    except Exception as e:
        logger.error(f"Failed to decrypt cookies: {e}")
        return ""


def delete_cookie() -> None:
    """Delete stored cookie."""
    if COOKIE_FILE.exists():
        COOKIE_FILE.unlink()
        logger.info("Cookies deleted from vault")


def is_configured() -> bool:
    """Check if cookie is configured."""
    if not COOKIE_FILE.exists():
        return False
    try:
        content = load_cookie()
        return bool(content and len(content) > 10)
    except Exception:
        return False


def has_content() -> bool:
    """Check if cookie vault has actual content."""
    return is_configured()
