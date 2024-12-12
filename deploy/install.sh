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
chmod +x $directory/bin/dbcli

# download other files needed
echo "Downloading run-migrations.sh script"
wget -q -O $directory/run-migrations.sh "https://raw.githubusercontent.com/felix-schott/jamsessions/refs/tags/$tag/backend/scripts/run-migrations.sh"

mkdir -p -m 755 $directory/init_db
echo "Downloading db-init scripts"
wget -q -O $directory/init_db/001_schema.sql "https://raw.githubusercontent.com/felix-schott/jamsessions/refs/tags/$tag/backend/internal/db/schema.sql"

wget -q -O $directory/init_db/002_roles.sh "https://raw.githubusercontent.com/felix-schott/jamsessions/refs/tags/$tag/backend/internal/db/scripts/add-roles.sh"

echo "Downloading production docker compose file"
wget -q -O $directory/docker-compose.yml "https://raw.githubusercontent.com/felix-schott/jamsessions/refs/tags/$tag/deploy/prod.docker-compose.yml"

echo "Downloading Caddyfile"
wget -q -O $directory/Caddyfile "https://raw.githubusercontent.com/felix-schott/jamsessions/refs/tags/$tag/deploy/Caddyfile"

echo "Downloading prometheus.yml"
wget -q -O $directory/prometheus.yml "https://raw.githubusercontent.com/felix-schott/jamsessions/refs/tags/$tag/deploy/prometheus.yml"

echo "Downloading migrations-alert.sh"
wget -q -O $directory/migrations-alert.sh "https://raw.githubusercontent.com/felix-schott/jamsessions/refs/tags/$tag/deploy/migrations-alert.sh"

[[ ! -f $directory/.env ]] && {
    echo ".env doesn't exist yet - creating file and default directories"
    touch $directory/.env
    echo "RELEASE_TAG=$tag" > $directory/.env
    echo "POSTGRES_DATA_DIR=$directory/postgres-data" >> $directory/.env
    echo "WEBSITE_HOST=:80" >> $directory/.env
    echo "POSTGRES_PASSWORD=replace-me" >> $directory/.env
    echo "READ_WRITE_PASSWORD=replace-me" >> $directory/.env
    echo "READ_ONLY_PASSWORD=replace-me" >> $directory/.env
    echo "PROD_UID=$UID" >> $directory/.env
    echo "PROD_GID=$UID" >> $directory/.env
    echo "LOCAL_DB_PORT=5432" >> $directory/.env
    echo "DB_URL=\"host=localhost port=\$LOCAL_DB_PORT user=read_write password=\$READ_WRITE_PASSWORD dbname=\$POSTGRES_DB sslmode=disable\"" >> $directory/.env
    echo "MIGRATIONS_DIRECTORY=$directory/migrations" >> $directory/.env
    echo "MIGRATIONS_ARCHIVE=$directory/migrations/archive" >> $directory/.env
    mkdir -p $directory/postgres-data
    mkdir -p $directory/migrations/suggestions
    mkdir -p $directory/migrations/archive
} || {
    echo ".env already exists - overwriting RELEASE_TAG"
    sed -i "s/RELEASE_TAG.*/RELEASE_TAG=$tag/g" $directory/.env
}

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
- LOCAL_DB_PORT (port to expose to the DB to - localhost only)
- DB_URL (postgres connection string - should include read-write credentials as this will be used for migrations)
- WEBSITE_HOST (the domain of the website to enable auto https, 0.0.0.0 for local bind)

If you wish, you can modify the docker-compose.yml file according to your needs. Note that running the default docker-compose won't work
if you're not the project owner, and you will have to build your own production docker images.

Then, you can start the application by running \`docker compose up -d\` in the directory $directory

# Managing the database
Whenever a user requests modification of the database (e.g. the addition of a new session), the application 
will write little bash scripts to $directory/migrations that make use the dbcli binary.

Review the script contents and execute run-migrations.sh to apply all changes.
EOF

echo "Installing alerting cron job"
set +eo pipefail
(crontab -l; echo "0 */2 * * * cd $directory && bash $directory/migrations-alert.sh $directory") | sort - | uniq - | crontab -;

echo "Finished installation process - please consult the generated README file for further instructions."