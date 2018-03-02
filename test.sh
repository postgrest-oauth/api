#!/bin/bash

cd ./example
docker-compose -f docker-compose.yml up -d --build
cd ../
newman run ./postgrest-oauth-server.postman_collection.json --bail --ignore-redirects --global-var host=localhost
