#!/bin/bash
echo "Building webapi project..."
docker build -t webapigo .
echo "Open integration folder"
cd integration

echo "Build integration testing project"
docker build -t webapitestgo .

echo "Go back main folder"
cd ..

echo "webapi are starting..."
docker-compose -f docker-compose.yml -f docker-compose.dev.yml up