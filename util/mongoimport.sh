#!/bin/sh

mongoimport \
	--db 'buenavida' \
	--collection 'products' \
	--file '/data/json/products.json' \
	--jsonArray \
	--uri "mongodb://$MONGO_USER:$MONGO_PASS@$MONGO_HOST:27017" \
	--authenticationDatabase 'admin'

mongoimport \
	--db 'buenavida' \
	--collection 'users' \
	--file '/data/json/users.json' \
	--jsonArray \
	--uri "mongodb://$MONGO_USER:$MONGO_PASS@$MONGO_HOST:27017" \
	--authenticationDatabase 'admin'
