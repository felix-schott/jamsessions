# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Calendar Versioning](https://calver.org/) (`YYYY.MM.MICRO-TAG`).

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