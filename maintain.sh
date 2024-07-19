#!/bin/bash

# set -x
ts=`TZ='Asia/shangh' date +%s -d '-3 month'` 
sql="DELETE from chain_transfer where ts <= $((ts))"

if [ $# -ne 4 ]; then
    echo "usage: ./manitain.sh  postgres localhost 5432 9527"
    exit
fi

dbname="sync_chain_"$4
echo psql -U $1 -h $2 -p $3 -d $dbname -c "$sql"
echo "$((ts)) is `TZ='Asia/shangh' date -d @$((ts))`"
psql -U $1 -h $2 -p $3 -d $dbname -c "$sql"