## ADDED Requirements

### Requirement: Cookie storage migration
The Gateway SHALL store cookies in `~/.xhs-gateway/cookies.enc` instead of `Spider_XHS/.cookies`.

#### Scenario: Migrate existing cookies on startup
- **WHEN** Gateway starts and `Spider_XHS/.cookies` exists but `~/.xhs-gateway/cookies.enc` does not
- **THEN** Gateway reads the old file, migrates to new location, deletes old file, logs "Cookies migrated"

#### Scenario: AES-256-GCM encrypted storage
- **WHEN** cookies are stored
- **THEN** they SHALL be encrypted using AES-256-GCM with a machine-specific key derived from hostname + username

#### Scenario: GET /api/v1/cookies returns content summary
- **WHEN** GET /api/v1/cookies is called
- **THEN** response includes `{"configured": true/false, "has_content": true/false, "domain": "xiaohongshu.com"}`

### Requirement: Cookie CRUD operations
The Gateway SHALL provide full cookie lifecycle management.

#### Scenario: Save cookies
- **WHEN** POST /api/v1/cookies with valid cookie string
- **THEN** Gateway encrypts and stores to `~/.xhs-gateway/cookies.enc`, returns success

#### Scenario: Delete cookies
- **WHEN** DELETE /api/v1/cookies
- **THEN** Gateway removes `cookies.enc`, returns success

#### Scenario: Cookie expiration handling
- **WHEN** Spider_XHS reports cookie as invalid (HTTP 401)
- **THEN** Gateway marks cookie as expired, returns `{"error": "cookie_expired"}` on next request
