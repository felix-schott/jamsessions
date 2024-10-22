#!/usr/bin/env bash

# script to bootstrap a deployment without local copy of the repo
# downloads all relevant files and runs the application

set -eo pipefail

[[ $1 == "" ]] && echo "Please provide the tag (package version) you want to install as the first positional argument." 1>&2 && exit 1;
tag=$1

[[ $2 == "" ]] && echo "Please provide the target directory of the installation as the second positional argument." 1>&2 && exit 1;
directory=$2

# download dbcli binary
echo "Downloading dbcli binary to ./bin/dbcli"
mkdir -p -m 755 $directory/bin
wget -q -O $directory/bin/dbcli "https://github.com/felix-schott/jamsessions/releases/download/$tag/dbcli"

# download run-migrations.sh script
echo "Downloading run-migrations.sh script"
wget -q -O $directory/run-migrations.sh "https://raw.githubusercontent.com/felix-schott/jamsessions/refs/tags/$tag/backend/scripts/run-migrations.sh"

echo "Downloading production docker compose file"
wget -q -O $directory/docker-compose.yml "https://raw.githubusercontent.com/felix-schott/jamsessions/refs/tags/$tag/deploy/prod.docker-compose.yml"

echo "Downloading Caddyfile"
wget -q -O $directory/Caddyfile "https://raw.githubusercontent.com/felix-schott/jamsessions/refs/tags/$tag/deploy/Caddyfile"

echo "Downloading prometheus.yml"
wget -q -O $directory/prometheus.yml "https://raw.githubusercontent.com/felix-schott/jamsessions/refs/tags/$tag/deploy/prometheus.yml"

echo "Creating readme file with instructions"
cat << EOF > $directory/README.md
# Starting the application
First, make sure there is a .env file present in $directory that contains the following variables:
- POSTGRES_DATA_DIR (local directory to persist the database to)
- POSTGRES_PASSWORD (password of the db superuser)
- READ_ONLY_PASSWORD (password for read-only db user)
- READ_WRITE_PASSWORD (password for rw db user)
- POSTGRES_DB (name of the database)
- PROD_UID (host uid that you want files in volumes to be owned by)
- PROD_GID (group id of image user)

If you wish, you can modify the docker-compose.yml file according to your needs. Note that running the default docker-compose won't work
if you're not the project owner, and you will have to build your own production docker images.

Then, you can start the application by running 'docker compose up -d' in the directory $directory

# Managing the database
Whenever a user requests modification of the database (e.g. the addition of a new session), the application 
will write little bash scripts to $directory/migrations that make use the dbcli binary.

Review the script contents and use the run-migrations.sh to apply all changes.
EOF

echo "Finished installation process - please consule the generate README file for further instructions."