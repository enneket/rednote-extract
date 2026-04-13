## ADDED Requirements

### Requirement: Cookie management
The system SHALL provide a UI for users to input and save their Xiaohongshu cookies. Cookies SHALL be stored locally in an encrypted file.

#### Scenario: Save cookies
- **WHEN** user pastes cookie string and clicks "Save"
- **THEN** system encrypts and stores cookie to local file and shows success message

#### Scenario: Load existing cookies
- **WHEN** user opens the client
- **THEN** system loads encrypted cookies from local file and displays "Cookie loaded"

#### Scenario: Clear cookies
- **WHEN** user clicks "Clear" button
- **THEN** system deletes stored cookies and resets cookie status to "Not configured"

### Requirement: Single note collection
The system SHALL allow users to input a single Xiaohongshu note URL and trigger data collection through the HTTP Gateway.

#### Scenario: Collect single note
- **WHEN** user enters a note URL and clicks "Collect"
- **THEN** system sends request to Gateway API and displays the collected note info (title, author, likes, etc.)

#### Scenario: Invalid note URL
- **WHEN** user enters an invalid note URL and clicks "Collect"
- **THEN** system shows error message "Invalid note URL format"

### Requirement: User notes collection
The system SHALL allow users to input a Xiaohongshu user profile URL and collect all notes from that user.

#### Scenario: Collect user notes
- **WHEN** user enters a user profile URL and clicks "Collect All"
- **THEN** system collects all notes and displays the count and list of collected items

### Requirement: Search and collect notes
The system SHALL allow users to input search keywords and collect matching notes.

#### Scenario: Search notes by keyword
- **WHEN** user enters keyword, selects sort type and note count, then clicks "Search"
- **THEN** system collects matching notes and displays results

#### Scenario: Search with filters
- **WHEN** user sets note type filter (video/image) and time filter
- **THEN** system applies filters to search results

### Requirement: Collection result display
The system SHALL display collected data in a structured format, showing note title, author, publish time, likes, comments, and bookmarks.

#### Scenario: Display collected notes
- **WHEN** notes are successfully collected
- **THEN** system displays results in a scrollable list with key info for each note

### Requirement: Export collected data
The system SHALL allow users to export collected notes to Excel file.

#### Scenario: Export to Excel
- **WHEN** user selects collected notes and clicks "Export"
- **THEN** system generates an Excel file and triggers download

### Requirement: Cross-platform packaging
The system SHALL be packaged as native executables for Windows (.exe), macOS (.dmg), and Linux (.AppImage) using Tauri bundler.

#### Scenario: Windows executable
- **WHEN** build command is run on Windows
- **THEN** a Windows .exe installer/package is generated

#### Scenario: macOS executable
- **WHEN** build command is run on macOS
- **THEN** a macOS .dmg file is generated

#### Scenario: Linux executable
- **WHEN** build command is run on Linux
- **THEN** a Linux .AppImage file is generated
