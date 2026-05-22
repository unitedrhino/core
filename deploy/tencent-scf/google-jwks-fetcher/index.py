"""腾讯云 SCF Google JWKS 获取器。

部署在境外地域，用于从 Google 官方 JWKS 地址获取公钥集合，并通过函数 URL
提供给国内生产服务缓存使用。
"""

import datetime
import hmac
import hashlib
import json
import os
import time
import urllib.error
import urllib.request


GOOGLE_JWKS_URL = "https://www.googleapis.com/oauth2/v3/certs"
DEFAULT_MAX_AGE_SECONDS = 21600
ACCESS_TOKEN_ENV = "YKHL_JWKS_ACCESS_TOKEN"
ACCESS_TOKEN_HEADER = "x-ykhl-jwks-token"

_cached_body = ""
_cached_headers = {}
_cache_expires_at = 0.0


def _parse_max_age(cache_control):
    """解析 Google 返回的 Cache-Control max-age。"""
    if not cache_control:
        return DEFAULT_MAX_AGE_SECONDS
    for part in cache_control.split(","):
        item = part.strip().lower()
        if item.startswith("max-age="):
            try:
                return max(60, int(item.split("=", 1)[1]))
            except ValueError:
                return DEFAULT_MAX_AGE_SECONDS
    return DEFAULT_MAX_AGE_SECONDS


def _fetch_google_jwks():
    """从 Google 官方地址获取并校验 JWKS 基本结构。"""
    req = urllib.request.Request(
        GOOGLE_JWKS_URL,
        headers={
            "Accept": "application/json",
            "User-Agent": "ykhl-google-jwks-fetcher/1.0",
        },
        method="GET",
    )
    with urllib.request.urlopen(req, timeout=10) as resp:
        raw = resp.read()
        jwks = json.loads(raw.decode("utf-8"))
        keys = jwks.get("keys")
        if not isinstance(keys, list) or len(keys) == 0:
            raise ValueError("google jwks keys is empty")
        for key in keys:
            if key.get("kty") != "RSA" or not key.get("kid") or not key.get("n") or not key.get("e"):
                raise ValueError("google jwks contains invalid key")
        canonical_jwks = json.dumps(jwks, ensure_ascii=False, sort_keys=True, separators=(",", ":")).encode("utf-8")
        now = datetime.datetime.utcnow().replace(microsecond=0).isoformat() + "Z"
        out = {
            "keys": keys,
            "_meta": {
                "issuer": "https://accounts.google.com",
                "source": GOOGLE_JWKS_URL,
                "fetchedAt": now,
                "googleCacheControl": resp.headers.get("Cache-Control", ""),
                "jwksSha256": hashlib.sha256(canonical_jwks).hexdigest(),
            },
        }
        body = json.dumps(out, ensure_ascii=False, sort_keys=True, separators=(",", ":"))
        max_age = _parse_max_age(resp.headers.get("Cache-Control", ""))
        headers = {
            "Content-Type": "application/json; charset=utf-8",
            "Cache-Control": "public, max-age=%d" % max_age,
            "X-Google-JWKS-SHA256": out["_meta"]["jwksSha256"],
        }
        return body, headers, time.time() + max_age


def _get_jwks_body(force_refresh=False):
    """读取缓存；缓存过期或强制刷新时重新请求 Google。"""
    global _cached_body, _cached_headers, _cache_expires_at
    now = time.time()
    if (not force_refresh) and _cached_body and now < _cache_expires_at:
        return _cached_body, _cached_headers
    body, headers, expires_at = _fetch_google_jwks()
    _cached_body = body
    _cached_headers = headers
    _cache_expires_at = expires_at
    return body, headers


def _event_headers(event):
    """从函数 URL 事件中读取请求头。"""
    if not isinstance(event, dict):
        return {}
    headers = event.get("headers") or event.get("Headers") or event.get("headerParameters") or {}
    if not isinstance(headers, dict):
        return {}
    return {str(k).lower(): str(v) for k, v in headers.items()}


def _is_http_event(event):
    """判断是否为函数 URL / HTTP 触发事件。"""
    if not isinstance(event, dict):
        return False
    return bool(event.get("headers") or event.get("Headers") or event.get("httpMethod") or event.get("requestContext"))


def _authorized(event):
    """校验生产后端请求函数 URL 时携带的共享密钥。"""
    if not _is_http_event(event):
        return True
    expected = os.environ.get(ACCESS_TOKEN_ENV, "")
    if expected == "":
        return False
    provided = _event_headers(event).get(ACCESS_TOKEN_HEADER, "")
    return hmac.compare_digest(provided, expected)


def main_handler(event, context):
    """main_handler 处理函数 URL、定时触发和手动 Invoke。"""
    force_refresh = False
    if isinstance(event, dict):
        force_refresh = event.get("type") == "timer" or event.get("forceRefresh") is True
    if not _authorized(event):
        return {
            "isBase64Encoded": False,
            "statusCode": 403,
            "headers": {"Content-Type": "application/json; charset=utf-8"},
            "body": json.dumps({"error": "forbidden"}, ensure_ascii=False),
        }
    try:
        body, headers = _get_jwks_body(force_refresh)
        return {
            "isBase64Encoded": False,
            "statusCode": 200,
            "headers": headers,
            "body": body,
        }
    except (urllib.error.URLError, TimeoutError, ValueError, json.JSONDecodeError) as err:
        return {
            "isBase64Encoded": False,
            "statusCode": 502,
            "headers": {"Content-Type": "application/json; charset=utf-8"},
            "body": json.dumps({"error": "fetch_google_jwks_failed", "message": str(err)}, ensure_ascii=False),
        }
