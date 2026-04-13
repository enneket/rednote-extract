## ADDED Requirements

### Requirement: Health check endpoint
The Gateway SHALL expose a GET /health endpoint that returns the service status.

#### Scenario: Health check
- **WHEN** client sends GET /health request
- **THEN** system returns `{"status": "ok", "spider_ready": true/false}`

### Requirement: Collect single note API
The Gateway SHALL expose POST /api/v1/notes/single that accepts a note URL and returns collected note data.

#### Scenario: Collect note via API
- **WHEN** POST /api/v1/notes/single is called with `{"note_url": "https://..."}`
- **THEN** system returns collected note data as JSON

#### Scenario: Missing note_url parameter
- **WHEN** POST /api/v1/notes/single is called without note_url
- **THEN** system returns 400 error with `{"error": "note_url is required"}`

### Requirement: Collect user notes API
The Gateway SHALL expose POST /api/v1/notes/user that accepts a user profile URL and returns all notes from that user.

#### Scenario: Collect user notes via API
- **WHEN** POST /api/v1/notes/user is called with `{"user_url": "https://..."}`
- **THEN** system returns list of notes as JSON array

### Requirement: Search notes API
The Gateway SHALL expose POST /api/v1/notes/search that accepts search parameters and returns matching notes.

#### Scenario: Search notes via API
- **WHEN** POST /api/v1/notes/search is called with `{"keyword": "榴莲", "count": 10}`
- **THEN** system returns matching notes as JSON array

#### Scenario: Search with all filters
- **WHEN** POST /api/v1/notes/search is called with full filter params
- **THEN** system applies all filters (sort_type, note_type, time, range) and returns filtered results

### Requirement: Cookie management API
The Gateway SHALL expose GET/POST /api/v1/cookies for managing cookies used in requests.

#### Scenario: Set cookies
- **WHEN** POST /api/v1/cookies is called with `{"cookies": "xsec_token=..."}`
- **THEN** system stores cookies for subsequent requests

#### Scenario: Get cookies
- **WHEN** GET /api/v1/cookies is called
- **THEN** system returns `{"configured": true/false}`

### Requirement: Error handling
The Gateway SHALL return structured error responses with appropriate HTTP status codes.

#### Scenario: Spider process error
- **WHEN** Spider_XHS subprocess fails
- **THEN** system returns 500 error with `{"error": "spider_error", "message": "..."}`

#### Scenario: Invalid URL format
- **WHEN** note_url or user_url is malformed
- **THEN** system returns 400 error with `{"error": "invalid_url", "message": "..."}`
