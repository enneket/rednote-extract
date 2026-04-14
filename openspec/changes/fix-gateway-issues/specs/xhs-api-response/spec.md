## ADDED Requirements

### Requirement: Unified note data structure
All note collection APIs SHALL return a unified note data structure with the following fields.

#### Scenario: Note structure fields
- **WHEN** any note API returns data
- **THEN** each note SHALL contain: `id`, `title`, `author`, `author_id`, `type`, `liked`, `collected`, `commented`, `shared`, `url`, `tags`, `time`, `ip_location`

### Requirement: Single note returns full data
The `/api/v1/notes/single` endpoint SHALL return full note information.

#### Scenario: Collect single note with full data
- **WHEN** POST /api/v1/notes/single with valid note_url
- **THEN** Gateway returns note object with all unified fields

### Requirement: User notes returns full data array
The `/api/v1/notes/user` endpoint SHALL return array of full note information.

#### Scenario: Collect user notes with full data
- **WHEN** POST /api/v1/notes/user with valid user_url
- **THEN** Gateway fetches all note URLs, then fetches full info for each, returns array of note objects

### Requirement: Search returns full data array
The `/api/v1/notes/search` endpoint SHALL return array of full note information.

#### Scenario: Search notes with full data
- **WHEN** POST /api/v1/notes/search with valid keyword and count
- **THEN** Gateway fetches matching note URLs, then fetches full info for each, returns array of note objects

### Requirement: Graceful degradation on partial failures
When collecting multiple notes, partial failures SHALL NOT abort the entire collection.

#### Scenario: Some notes fail during batch collection
- **WHEN** batch collecting notes and some fail
- **THEN** Gateway returns successfully collected notes with `partial: true` flag and `errors` array listing failed note IDs
