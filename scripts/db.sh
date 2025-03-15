#!/bin/bash

path="$(pwd)/scripts"
main_dir=$(dirname "$path")

get_conn_string() {
    . $path/loadEnv.sh && echo "postgresql://$POSTGRES_USER:$POSTGRES_PASSWORD@$POSTGRES_HOST:$POSTGRES_PORT/shared?sslmode=$POSTGRES_SSL_MODE"
}

run_psql() {
    conn_string=$(get_conn_string)
    psql $conn_string
}

create_pgpass() {
    . $path/loadEnv.sh && echo "$POSTGRES_HOST:$POSTGRES_PORT:shared:$POSTGRES_USER:$POSTGRES_PASSWORD" > $main_dir/.pgpass
    chmod 600 $main_dir/.pgpass
}

main() {
    create_pgpass "$@"
    run_psql
    exit 0
}

main
