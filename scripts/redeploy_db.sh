#!/bin/bash

docker container stop crype-db
docker container rm crype-db
docker volume rm crype_data

cd .. && docker compose up -d