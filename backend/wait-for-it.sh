#!/usr/bin/env bash

set -e

host="$1"
shift
cmd="$@"

until nc -z -v -w30 $host 3307; do
  >&2 echo "MySQL is unavailable - sleeping"
  sleep 1
done

>&2 echo "MySQL is up - executing command"
exec $cmd
