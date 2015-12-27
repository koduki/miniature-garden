#!/bin/bash

WORKDIR=/home/koduki/eigo-de-news/


cd ${WORKDIR}

# update
git pull
docker-compose -f docker-compose-prod.yml pull

# run
docker-compose -f docker-compose-prod.yml run app ruby /app/src/main.rb
