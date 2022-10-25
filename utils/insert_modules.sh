#!/bin/sh
set -e

until pg_isready -h $1 -U postgres
do
  echo "Waiting for postgres"
  sleep 4;
done

#psql -U postgres -h localhost -p 5432 -U postgres  < cpu_usage.sql
PGPASSWORD=rdsdb psql -U postgres -d homework -h localhost -p 5432  -c "\COPY cpu_usage FROM cpu_usage.csv CSV HEADER"
