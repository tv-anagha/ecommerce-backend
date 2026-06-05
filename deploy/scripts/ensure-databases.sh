#!/bin/sh
set -e

PGHOST="${PGHOST:-postgres}"
PGUSER="${PGUSER:-postgres}"
PGPASSWORD="${PGPASSWORD:-postgres}"
export PGPASSWORD

for db in cart_db order_db; do
  exists="$(psql -h "$PGHOST" -U "$PGUSER" -tc "SELECT 1 FROM pg_database WHERE datname = '${db}'")"
  if [ "$(echo "$exists" | tr -d '[:space:]')" != "1" ]; then
    echo "creating database ${db}"
    psql -h "$PGHOST" -U "$PGUSER" -c "CREATE DATABASE ${db};"
  else
    echo "database ${db} already exists"
  fi
done
