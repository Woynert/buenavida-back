#!/bin/bash

printmsg()
{
	printf "\
=====================================\n\
$1\n\
=====================================\n"
}

sudo rm -r ./data/mongodb/
sudo rm -r ./data/postgres/

mkdir ./data/mongodb
mkdir ./data/postgres


printmsg "Setting up MongoDB"

docker-compose -f ./util/import-mongo.yml up \
	--abort-on-container-exit \
	--exit-code-from mongo-script

if [ $? -eq 0 ]; then
	printmsg "MongoDB setup: Finished successfully"
else
	printmsg "MongoDB setup: Finished with error"
	exit 1
fi


printmsg "Setting up Postgres"

docker-compose -f ./util/setup-postgres.yml up \
	--abort-on-container-exit \
	--exit-code-from postgres-script

if [ $? -eq 0 ]; then
	printmsg "Postgres setup: Finished Successfully"
else
	printmsg "Postgres setup: Finished with error"
fi


