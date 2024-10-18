#!/usr/bin/env bash

# applies all changes in $MIGRATIONS_DIRECTORY and moves the scripts to 
# the archive ($MIGRATIONS_ARCHIVE) afterwards.
# use -y flag for non-interactive mode

# if you're not running this script as part of the production setup (see deploy/install.sh)
# make sure the dbcli binary is on your PATH

set -euo pipefail

# load .env file if it exists
[[ ! -f .env ]] || export $(grep -v '^#' .env | xargs)

# if there is a local bin directory (normal production setup), add to PATH
[[ -d "$PWD/bin" ]] && export PATH=$PATH:$PWD/bin

[[ $MIGRATIONS_DIRECTORY == "" ]] && echo "Please provide the environment variable 'MIGRATIONS_DIRECTORY'" 1>&2 && exit 1;

if [ -z "$( ls -Ap $MIGRATIONS_DIRECTORY | grep -v / )" ] 
then # list all files in the directory (make ls append / to directories, then filter)
   echo "The directory $MIGRATIONS_DIRECTORY is empty, no migrations to run" 1>&2;
else
  if [ $1 == "-y" ]
  then
    echo "Running in non-interactive mode." 1>&2;
    choice="y"
  else
    # wait for confirmation from user
    read -p "Apply all changes in $MIGRATIONS_DIRECTORY (y/n)?" choice
  fi

  [[ $MIGRATIONS_ARCHIVE == "" ]] && echo "Please provide the environment variable 'MIGRATIONS_ARCHIVE'" 1>&2 && exit 1;

  mkdir -p -m 755 $MIGRATIONS_ARCHIVE

  out=""

  # if yes, run all files in migrations directory
  case "$choice" in 
    y|Y ) 
      for file in $MIGRATIONS_DIRECTORY/*
      do 
        if [[ -f $file ]]
        then
          echo "Running $file" 1>&2;
          out="$out $(bash $file)";
          echo "Moving file to archive $MIGRATIONS_ARCHIVE/" 1>&2 && mv $file $MIGRATIONS_ARCHIVE/;
          echo -e "\n" 1>&2;
        fi
      done
      ;;
    n|N ) 
      echo "Not applying any changes" 1>&2;
      ;;
    * ) 
      echo "invalid choice $choice'" 1>&2;
      ;;
  esac
fi

echo $out # return stdout of all commands