#!/bin/sh

host="$1"
shift
cmd="$@"

until nc -z "$host" 5432; do
  echo "Waiting for $host..."
  sleep 2
done

exec $cmd
