## ADDED Requirements

### Requirement: Spider_XHS path auto-discovery
The Gateway SHALL automatically discover the Spider_XHS installation path on startup using a priority order.

#### Scenario: Environment variable takes precedence
- **WHEN** environment variable `SPIDER_XHS_PATH` is set
- **THEN** Gateway uses that path and logs "Using Spider_XHS from SPIDER_XHS env"

#### Scenario: Pip package installation detected
- **WHEN** Spider_XHS is installed as a pip package (`pip install spider-xhs`)
- **THEN** Gateway locates it via `importlib.util.find_spec` and logs "Using Spider_XHS from pip package"

#### Scenario: Sibling directory for development
- **WHEN** no env var and no pip package, but `../Spider_XHS` directory exists
- **THEN** Gateway uses the sibling path and logs "Using Spider_XHS from sibling directory"

#### Scenario: Spider_XHS not found
- **WHEN** none of the above methods succeed
- **THEN** Gateway raises `StartupError: Spider_XHS not found` with instructions to set `SPIDER_XHS_PATH`

### Requirement: Gateway configuration file
The Gateway SHALL read/write configuration from `~/.xhs-gateway/config.yaml`.

#### Scenario: Config directory creation
- **WHEN** Gateway starts and `~/.xhs-gateway/` does not exist
- **THEN** Gateway creates the directory and default config file

#### Scenario: Config stores Spider_XHS path override
- **WHEN** user sets Spider_XHS path via admin API
- **THEN** Gateway persists it to `~/.xhs-gateway/config.yaml`
