## 1. Project Setup

- [ ] 1.1 Create `apps/gateway/` directory with FastAPI project structure
- [ ] 1.2 Add dependencies: fastapi, uvicorn, pydantic, python-dotenv
- [ ] 1.3 Create `apps/client/` directory with Tauri + React + TypeScript project
- [ ] 1.4 Initialize Tauri with `npm create tauri-app`
- [ ] 1.15 Configure Tauri for cross-platform build (Windows/macOS/Linux)

## 2. HTTP Gateway Implementation

- [ ] 2.1 Implement `/health` endpoint
- [ ] 2.2 Implement `POST /api/v1/notes/single` endpoint
- [ ] 2.3 Implement `POST /api/v1/notes/user` endpoint
- [ ] 2.4 Implement `POST /api/v1/notes/search` endpoint
- [ ] 2.5 Implement `GET/POST /api/v1/cookies` endpoint
- [ ] 2.6 Integrate Spider_XHS as Python subprocess in Gateway
- [ ] 2.7 Add structured error handling and logging

## 3. Desktop Client UI Implementation

- [ ] 3.1 Create main layout with navigation (Cookie / Single Note / User / Search)
- [ ] 3.2 Implement Cookie management page (input, save, load, clear)
- [ ] 3.3 Implement Single Note collection page
- [ ] 3.4 Implement User notes collection page
- [ ] 3.5 Implement Search notes page with filters (sort, type, time)
- [ ] 3.6 Implement results display component
- [ ] 3.7 Implement Excel export functionality
- [ ] 3.8 Add loading states and error handling in UI

## 4. Integration & Testing

- [ ] 4.1 Connect Tauri client to Gateway via HTTP (with localhost fallback for dev)
- [ ] 4.2 Test full flow: cookie config → single note → user notes → search
- [ ] 4.3 Cross-platform build test (at least Windows exe)

## 5. Packaging & Distribution

- [ ] 5.1 Configure Tauri bundler for Windows (.exe), macOS (.dmg), Linux (.AppImage)
- [ ] 5.2 Build production executables
- [ ] 5.3 Test that built executables run standalone
