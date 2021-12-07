#!/bin/bash
cd manifests/postgres/local
docker stop test-postgres
docker rm test-postgres
sleep 5
docker build -t test-postgres .
docker run --rm --name test-postgres -d -p 5432:5432 test-postgres