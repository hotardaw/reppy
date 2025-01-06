#!/bin/bash

docker-compose up -d --build
sleep 5

# Open tabs in chrome
# open http://localhost:8080    # frontend
# open http://localhost:8081    # backend
open http://localhost:8083    # pgAdmin