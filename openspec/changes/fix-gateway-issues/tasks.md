## 1. Gateway Path Discovery & Config

- [x] 1.1 Create `apps/gateway/app/core/config.py` with `discover_spider_xhs()` using env/pip/sibling priority
- [x] 1.2 Create `~/.xhs-gateway/` directory and `config.yaml` on startup
- [x] 1.3 Add `openpyxl` and `cryptography` to `requirements.txt`
- [x] 1.4 Update `main.py` lifespan to use new config discovery, fail fast if Spider_XHS not found

## 2. Cookie Vault

- [x] 2.1 Create `apps/gateway/app/core/cookie_vault.py` with AES-256-GCM encrypt/decrypt
- [x] 2.2 Implement cookie migration from `Spider_XHS/.cookies` to `~/.xhs-gateway/cookies.enc`
- [x] 2.3 Update `apps/gateway/app/api/cookies.py`: GET returns content summary, DELETE supported, migrate legacy cookie on startup
- [x] 2.4 Add DELETE /api/v1/cookies endpoint

## 3. Unified API Response

- [x] 3.1 Define `NoteData` Pydantic model with unified fields: id, title, author, author_id, type, liked, collected, commented, shared, url, tags, time, ip_location
- [x] 3.2 Update `apps/gateway/app/api/notes.py`: `/single` maps handle_note_info output to NoteData
- [x] 3.3 Update `/user` endpoint: fetch all note URLs then batch fetch full info for each, return NoteData array
- [x] 3.4 Update `/search` endpoint: same batch fetch pattern, return NoteData array
- [x] 3.5 Add `partial: bool` and `errors: list` to response for graceful degradation
- [x] 3.6 Update `apps/client/src/api/client.ts`: update TypeScript types to match NoteData

## 4. Excel Export

- [x] 4.1 Create `apps/gateway/app/api/export.py` with POST /api/v1/notes/export endpoint
- [x] 4.2 Integrate Spider_XHS `save_to_xlsx` via openpyxl for Excel generation
- [x] 4.3 Return file as StreamingResponse with Content-Disposition attachment header
- [x] 4.4 Support `filename` query parameter for custom export name

## 5. Verification

- [x] 5.1 Gateway starts and discovers Spider_XHS path (log output verification)
- [x] 5.2 Cookie save/load roundtrip with encryption verified
- [x] 5.3 `/single` returns full NoteData fields
- [x] 5.4 `/user` and `/search` return NoteData arrays (may need real cookies/test data)
- [x] 5.5 `/export` generates valid .xlsx download
