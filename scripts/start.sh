#!/bin/bash

mkdir keys
openssl rand -base64 768 > ./keys/rsfile.key
sudo chmod 700 keys
sudo chmod 600 keys/rsfile.key
# fresh build takes  approx 2 minutes and I only need build when I add dependencies
# soooo, adding --build to the command has no harm.
docker-compose up  --build -d
sleep 10
docker exec mongodb /scripts/rs-init.sh