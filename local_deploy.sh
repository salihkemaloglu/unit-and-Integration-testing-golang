#!/bin/bash
echo "Building webapi project..."
docker build -t webapigo .

echo "webapi are starting..."
docker-compose -f docker-compose.yml -f docker-compose.dev.yml up