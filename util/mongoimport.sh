#!/bin/sh

#MONGOHOST="localhost"

# start daemon
#mongod &

## import file
#sleep 15

mongoimport \
	--db 'buenavida' \
	--collection 'products' \
	--file '/data/products.json' \
	--jsonArray \
	--uri "mongodb://root:example@$MONGOHOST:27017" \
	--authenticationDatabase 'admin'

#sleep infinity
