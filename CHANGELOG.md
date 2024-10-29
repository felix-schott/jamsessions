# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Calendar Versioning](https://calver.org/) (`YYYY.MM.MICRO-TAG`).

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