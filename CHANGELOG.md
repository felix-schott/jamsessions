# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Calendar Versioning](https://calver.org/) (`YYYY.MM.MICRO-TAG`).

## [v2024.12.0] - 2024-12-02

### Changed

- Changed path structure from `/{sessionId}` to `/{venueName}-{venueId}/{sessionName}-{sessionId}` (PR [#94](https://github.com/felix-schott/jamsessions/pull/94))

### Fixed

- Fixed small bug on frontend that resulted in error messages after hitting the CLear button of the date input (PR [#90](https://github.com/felix-schott/jamsessions/pull/90))

### Added

- Added analytics using Plausible (PR [#100](https://github.com/felix-schott/jamsessions/pull/100))
- Added a Telegram Bot API integration for alerts when there are pending data migrations (PR [#97](https://github.com/felix-schott/jamsessions/pull/97))
- Added a venue page that lists all jam sessions happening at a certain venue (PRs [#94](https://github.com/felix-schott/jamsessions/pull/94), [#99](https://github.com/felix-schott/jamsessions/pull/99))
- Added genres 'Hip-Hop' and 'RnB' (PR [#86](https://github.com/felix-schott/jamsessions/pull/86))
- Added validation for string enums like `Genre` and `Backline` using custom `MarshalJSON` methods (PR [#86](https://github.com/felix-schott/jamsessions/pull/86))

## [v2024.11.4-beta] - 2024-11-19

### Fixed

- Optimised UX when switching between list and map view (PR [#79](https://github.com/felix-schott/jamsessions/pull/79))

### Added

- Added functionality for optional rating to be submitted and stored alongside a comment (PR [#83](https://github.com/felix-schott/jamsessions/pull/83))

### Changed

- Changed content of InfoPopup, added link to github repo (PR [#80](https://github.com/felix-schott/jamsessions/pull/80))
- Changed order in which venues are displayed in `AddSessionPopup` to alphabetical (PR [#82](https://github.com/felix-schott/jamsessions/pull/82))

## [v2024.11.3-beta] - 2024-11-15

### Fixed

- Fixed bug in AddSessionPopup, where the venue foreign key was not added to the payload for sessions associated with existing venues (PR [#76](https://github.com/felix-schott/jamsessions/pull/76))
- Fixed bug in Postgres function `sessions_on_date()` that resulted in `NthOfMonth` sessions being incorrectly included in query results ([#72](https://github.com/felix-schott/jamsessions/pull/72))

### Added

- Added `Fortnightly` interval (PR [#75](https://github.com/felix-schott/jamsessions/pull/75))

## [v2024.11.2-beta] - 2024-11-13

### Added

- Added rate limiter to geocoding module (PR [#66](https://github.com/felix-schott/jamsessions/pull/66))

## [v2024.11.1-beta] - 2024-11-13

### Fixed

- Fixed bug in AddSessionPopup - time component of the start datetime was being ignored (PR [#62](https://github.com/felix-schott/jamsessions/pull/62))

## [v2024.11.0-beta] - 2024-11-13

### Added

- Added list view alongside map view (PR [#56](https://github.com/felix-schott/jamsessions/pull/56))
- Added date range filter to API (PR [#53](https://github.com/felix-schott/jamsessions/pull/53), [#55](https://github.com/felix-schott/jamsessions/pull/55))
- Added support for irregular sessions (PR [#57](https://github.com/felix-schott/jamsessions/pull/57))

## [v2024.10.4-beta] - 2024-10-30

### Fixed

- Several small bug fixes in frontend (PR [#38](https://github.com/felix-schott/jamsessions/pull/38)), backend (PRs [#37](https://github.com/felix-schott/jamsessions/pull/37), [#40](https://github.com/felix-schott/jamsessions/pull/40)) and deployment code (PR [#39](https://github.com/felix-schott/jamsessions/pull/39))

## [v2024.10.3-beta] - 2024-10-29

### Added

- Store author and time alongside comment content in separate DB table; implement UI (PR [#34](https://github.com/felix-schott/jamsessions/pull/34))
- Add rating functionality with new DB table/API routes and UI changes (PR [#34](https://github.com/felix-schott/jamsessions/pull/34))
- Add UI to change backline and genre information in `EditSessionPopup.svelte` (PR [#35](https://github.com/felix-schott/jamsessions/issues/35))

### Changed

- Upgrade frontend code to Svelte 5 (PR [#36](https://github.com/felix-schott/jamsessions/pull/34))

### Fixed

- Fixed a bug that occurred with an empty database (PR [#28](https://github.com/felix-schott/jamsessions/pull/28))
- Propagate release tag from `install.sh` to the docker image version in `docker-compose.yml` (PR [#33](https://github.com/felix-schott/jamsessions/pull/33))

## [v2024.10.2-beta] - 2024-10-24

### Added

- Add prometheus container to docker compose

### Changed

- Switch from nginx to caddy for easy prometheus integration and https

## [v2024.10.1-beta] - 2024-10-18

### Added

- Added local deployment setup and `deploy/install.sh` script.

### Changed

- Small modifications to `scripts/run-migrations.sh` - load .env file automatically if it exists.

## [v2024.10.0-alpha] - 2024-10-17

### Added

- First more or less deployable version of the software.