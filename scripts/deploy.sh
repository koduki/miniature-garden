#!/bin/bash

cd /home/koduki/eigo-de-news/


git pull

docker-compose -f docker-compose-prod.yml pull
docker-compose -f docker-compose-prod.yml stop
docker-compose -f docker-compose-prod.yml up -d
docker-compose -f docker-compose-prod.yml ps
