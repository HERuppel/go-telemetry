#!/bin/bash

mongo -u $MONGO_INITDB_ROOT_USERNAME -p $MONGO_INITDB_ROOT_PASSWORD <<EOF
use telemetry;                
db.createCollection("events")
print("Database successfully created!");
EOF