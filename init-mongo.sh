#!/bin/bash

mongo -u $MONGO_INITDB_ROOT_USERNAME -p $MONGO_INITDB_ROOT_PASSWORD <<EOF
use $MONGO_DB_NAME;                
db.createCollection($MONGO_DB_COLLECTION)
db.$MONGO_DB_COLLECTION.createIndex({ "timestamp": 1 })
db.$MONGO_DB_COLLECTION.createIndex({ "type": 1 })
db.createCollection($MONGO_DB_METRICS_COLLECTION)
print("Database successfully created!");
EOF