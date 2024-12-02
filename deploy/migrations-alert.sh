#!/usr/bin/env bash

set -eo pipefail

[[ $1 == "" ]] && echo "Please provide the target directory of the installation as a positional argument." 1>&2 && exit 1;
directory=$1/migrations

[[ $TELEGRAM_TOKEN == "" ]] && echo "Please provide the environment variable TELEGRAM_TOKEN." 1>&2 && exit 1;
[[ $TELEGRAM_CHAT_ID == "" ]] && echo "Please provide the environment variable TELEGRAM_CHAT_ID." 1>&2 && exit 1;

send_alert () {
  curl -X POST https://api.telegram.org/bot${TELEGRAM_TOKEN}/sendMessage -d '{"chat_id":'${TELEGRAM_CHAT_ID}',"text":"PENDING MIGRATIONS"}' -H 'Content-Type: application/json'
}

if [[ -d $directory ]]; then
    ls $directory/*.sh &> /dev/null && send_alert &> /dev/null
else
    echo "Directory $directory does not exist" 1>&2 && exit 1
fi