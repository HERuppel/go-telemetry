#!/bin/bash

mongo <<EOF
use telemetry;                
db.createCollection("events")
print("Database successfully created!");
EOF