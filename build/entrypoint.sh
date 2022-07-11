#!/bin/sh
./wait-for-it/wait-for-it.sh database:5432 -- make migrate
exec "$@"