#!/bin/sh

export PGPASSWORD=$POSTGRES_PASSWORD

# wait till postgres server is up
sleep 10

# execute sql scripts
psql -U $POSTGRES_USER -h $POSTGRES_HOST -p 5432 -d $POSTGRES_DB -f /opt/sql/tables.sql
psql -U $POSTGRES_USER -h $POSTGRES_HOST -p 5432 -d $POSTGRES_DB -f /opt/sql/procedures.sql
psql -U $POSTGRES_USER -h $POSTGRES_HOST -p 5432 -d $POSTGRES_DB -f /opt/sql/sample-data.sql
