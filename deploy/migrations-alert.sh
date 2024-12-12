#!/usr/bin/env bash

set -eo pipefail

# load .env file if it exists
[[ ! -f .env ]] || set -a && source .env && set +a

[[ $1 == "" ]] && echo "Please provide the target directory of the installation as a positional argument." 1>&2 && exit 1;
directory=$1

[[ $TELEGRAM_TOKEN == "" ]] && echo "Please provide the environment variable TELEGRAM_TOKEN." 1>&2 && exit 1;
[[ $TELEGRAM_CHAT_ID == "" ]] && echo "Please provide the environment variable TELEGRAM_CHAT_ID." 1>&2 && exit 1;

send_alert () {
  type=$1
  echo "Sending $type alert"
  curl -X POST https://api.telegram.org/bot${TELEGRAM_TOKEN}/sendMessage -d '{"chat_id":'${TELEGRAM_CHAT_ID}',"text":"PENDING '$type'"}' -H 'Content-Type: application/json'
}

echo "Inspecting $directory/migrations directory"
if [[ -d "$directory/migrations" ]]; then
    ls $directory/migrations/*.sh &> /dev/null && send_alert "MIGRATIONS" &> /dev/null
else
    echo "Directory $directory/migrations does not exist" 1>&2 && exit 1
fi

echo "Inspecting $directory/migrations/suggestions directory"
if [[ -d "$directory/migrations/suggestions" ]]; then
    ls $directory/migrations/suggestions/*.sh &> /dev/null && send_alert "SUGGESTIONS" &> /dev/null
else
    echo "Directory $directory/migrations/suggestions does not exist" 1>&2 && exit 1
fi