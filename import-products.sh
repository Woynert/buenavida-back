#!/bin/sh

sudo rm -r ./data/mongodb/
sudo rm -r ./data/postgres/

mkdir ./data/mongodb
mkdir ./data/postgres

docker-compose -f ./util/import-products.yml up
