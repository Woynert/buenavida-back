#!/bin/sh

# import products.json

mongoimport \
	--db 'buenavida' \
	--collection 'products' \
	--file '/data/json/products.json' \
	--jsonArray \
	--uri "mongodb://$MONGO_USER:$MONGO_PASS@$MONGO_HOST:27017" \
	--authenticationDatabase 'admin'

# import duplicates

mongoimport \
	--db 'buenavida' \
	--collection 'products' \
	--file '/data/json/productsdup.json' \
	--jsonArray \
	--uri "mongodb://$MONGO_USER:$MONGO_PASS@$MONGO_HOST:27017" \
	--authenticationDatabase 'admin'

# import users.json

mongoimport \
	--db 'buenavida' \
	--collection 'users' \
	--file '/data/json/users.json' \
	--jsonArray \
	--uri "mongodb://$MONGO_USER:$MONGO_PASS@$MONGO_HOST:27017" \
	--authenticationDatabase 'admin'
