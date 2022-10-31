#!/bin/sh

sudo rm -r ./data/mongodb/
sudo rm -r ./data/postgres/

mkdir ./data/mongodb
mkdir ./data/postgres

docker-compose -f ./util/import-products.yml up \
	--abort-on-container-exit \
	--exit-code-from mongo-script

if [ $? -eq 0 ]; then
	printf "\
===========================\n\
Finished Successfully\n\
===========================\n"
else
	printf "\
===========================\n\
Finished with error\n\
===========================\n"
fi


