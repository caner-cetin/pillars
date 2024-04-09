#!/bin/bash
# these files are supplied by .env to docker-compose and to here
mongosh -u $MONGO_INITDB_ROOT_USERNAME -p $MONGO_INITDB_ROOT_PASSWORD <<EOF
var config = {
    "_id": "rs0",
    "version": 1,
    "members": [
        {
            "_id": 1,
            "host": "mongodb:27017",
            "priority": 1
        },
    ]
};
rs.initiate(config, { force: true });
rs.status();
EOF
