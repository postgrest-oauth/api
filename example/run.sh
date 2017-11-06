#!/bin/bash

docker-compose -f docker-compose.yml stop
docker-compose -f docker-compose.yml rm -f
docker volume ls -q | xargs docker volume rm
docker-compose -f docker-compose.yml up --build    
