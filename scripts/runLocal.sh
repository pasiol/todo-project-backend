#!/bin/bash
cd manifests/postgres/local
docker stop test-postgres
docker rm test-postgres
docker build -t test-postgres .
docker run --rm --name test-postgres -d -p 5432:5432 test-postgres
cd ../..
docker ps | grep test-postgres
docker run -d -p 8888:8888 -e APP_PORT=8888 -e ALLOWED_ORIGINS=http://localhost:3000 -e POSTGRES_USER=tester -e POSTGRES_PASSWORD=testing -e POSTGRES_DB=todos -e POSTGRES_HOST=localhost pasiol/todo-project-backend