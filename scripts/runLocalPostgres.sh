#!/bin/bash
cd manifests/postgres/local
docker stop test-postgres
docker build -t test-postgres .
docker run --rm --name test-postgres -d -p 5432:5432 test-postgres
cd ../..
docker ps | grep test-postgres