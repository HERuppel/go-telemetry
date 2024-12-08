#!/bin/bash

mongo -u $MONGO_INITDB_ROOT_USERNAME -p $MONGO_INITDB_ROOT_PASSWORD <<EOF
use $MONGO_DB_NAME;                
db.createCollection($MONGO_DB_COLLECTION)
db.createCollection($MONGO_DB_METRICS_COLLECTION)
print("Database successfully created!");
EOF