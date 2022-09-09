#!/bin/bash

set -e

if [[ $MUST_MIGRATE ]]
then
    echo running migrations...
    /opt/tg/bin/migrate -path /opt/tg/migrations -database ${DB_DSN} up
fi

exec "$@"