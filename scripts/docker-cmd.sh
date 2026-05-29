#!/bin/bash
set -e

export PGDATA=/var/lib/postgresql/data

# Start PostgreSQL in background
echo "Starting PostgreSQL..."
su postgres -c "pg_ctl start -D $PGDATA -w -t 30"

echo "PostgreSQL is ready!"

# Start Go server in foreground
echo "Starting consumable-tracker..."
exec /usr/local/bin/consumable-tracker
