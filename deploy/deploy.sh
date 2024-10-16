#!/usr/bin/env bash

tag=$1
[[ $tag == "" ]] && echo "Please provide a tag (package version) as a positional argument." 1>&2 && exit 1;

# needed
# - bin/dbcli (github pipeline artifact? or build locally)
# - docker image server (pull from ghcr - pipeline release)
# - docker image frontend (pull from ghcr - pipeline release)

# then run prod.docker-compose.yml

# also copy bin/run-migrations.sh (maybe makefile?)