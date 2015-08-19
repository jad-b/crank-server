#!/bin/bash -eu
#
# torque-ops
# ==========
#
# Provides some helpful command-line calls for dealing with Torque.
#
###############################################################################

# Variables
PSQL_VERSION=9.4.4

###############################################################################
# Save the postgres DB as a template file in current directory
#
# Args
#   1) Name of the DB container
###############################################################################
dump_db(){
    docker run \
        --rm \
        -it \
        -v $(pwd):/backup \
        --link "$1":db \
        "postgres:$PSQL_VERSION" \
        sh -c 'pg_dump -h db -U torque -W -d torque --create > /backup/torque.template.pg'
}

case $1 in
    savedb)
        dump_db "$2"
        ;;
esac
