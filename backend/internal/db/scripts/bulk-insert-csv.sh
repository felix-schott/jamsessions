#!/usr/bin/env bash

[[ $1 == "" ]] && {
    echo "Please provide the port as the first positional arg" 1>&2 && exit
} || {
    port=$1
}

[[ $2 == "" ]] && {
    echo "Please provide the csv filepath as the second positional arg" 1>&2 && exit
} || {
    file=$2
}

psql \
  -h localhost -p $port -d london_jam_sessions -U read_write \
  -c "\copy mytable (column1, column2)  from '$file' with delimiter as ','"