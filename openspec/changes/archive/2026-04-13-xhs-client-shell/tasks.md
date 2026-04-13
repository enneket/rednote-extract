## 1. Project Setup

- [x] 1.1 Create `apps/gateway/` directory with FastAPI project structure
- [x] 1.2 Add dependencies: fastapi, uvicorn, pydantic, python-dotenv
- [x] 1.3 Create `apps/client/` directory with Tauri + React + TypeScript project
- [x] 1.4 Initialize Tauri (manual scaffold since CLI requires tty)
- [x] 1.15 Configure Tauri for cross-platform build (Windows/macOS/Linux)

## 2. HTTP Gateway Implementation

- [x] 2.1 Implement `/health` endpoint
- [x] 2.2 Implement `POST /api/v1/notes/single` endpoint
- [x] 2.3 Implement `POST /api/v1/notes/user` endpoint
- [x] 2.4 Implement `POST /api/v1/notes/search` endpoint
- [x] 2.5 Implement `GET/POST /api/v1/cookies` endpoint
- [x] 2.6 Integrate Spider_XHS via direct Python import in Gateway
- [x] 2.7 Add structured error handling and logging

## 3. Desktop Client UI Implementation

- [x] 3.1 Create main layout with navigation (Cookie / Single Note / User / Search)
- [x] 3.2 Implement Cookie management page (input, save, load, clear)
- [x] 3.3 Implement Single Note collection page
- [x] 3.4 Implement User notes collection page
- [x] 3.5 Implement Search notes page with filters (sort, type, time)
- [x] 3.6 Implement results display component
- [ ] 3.7 Implement Excel export functionality (需要 Spider_XHS 的 save_to_xlsx，Gateway 目前只返回 URL 列表)
- [x] 3.8 Add loading states and error handling in UI

## 4. Integration & Testing

- [x] 4.1 Connect Tauri client to Gateway via HTTP (with localhost fallback for dev)
- [ ] 4.2 Test full flow: cookie config → single note → user notes → search
- [ ] 4.3 Cross-platform build test (at least Windows exe)

## 5. Packaging & Distribution

- [x] 5.1 Configure Tauri bundler for Windows (.exe), macOS (.dmg), Linux (.AppImage/.deb/.rpm)
- [x] 5.2 Build production executables (binary + .deb + .rpm built; AppImage failed due to linuxdeploy download)
- [ ] 5.3 Test that built executables run standalone
