## Overview
> IMPORTANT: This is a WIP!

This repository contains code for a yet-to-be-published web application cataloguing and presenting public jam sessions in London. 

The `backend` directory contains the Golang codebase for a REST API server alongside a CLI tool for database management. The `backend/internal/db` module contains database schema and `sqlc` models. Please refer to the `Makefile` to run tests and build processes. 

The `frontend` directory contains the Svelte codebase for the web interface of the application (MPA). Central to the UI is a map that shows the locations of sessions. The website allows users to filter jam sessions by date, genre and backline provided by the venue. It further gives user the ability to suggest changes to individual session entries or add new ones. All these suggestions end up on the server as database migrations that manually need to be approved of and run by the admin. The custom migration system makes use of aforementioned `dbcli` tool.

The `deploy` directory contains an `install.sh` script to bootstrap a local installation (with docker compose).

There are Github Actions set up to 
1. run the test suite as part of a PR
2. build and release frontend and backend docker images after a tag is created
3. build and release the `dbcli` binary after a tag is created

`CONTRIBUTING.md` with instructions on how to contribute is following!