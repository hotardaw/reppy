#!/bin/bash

IMAGES=(
  "postgres:14"
  "dpage/pgadmin4"
  "node:latest"
  "golang:latest"
)

for IMAGE in "${IMAGES[@]}"
do
  echo "Pulling $IMAGE..."
  docker pull $IMAGE
done