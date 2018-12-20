#!/bin/bash

BLUE='\033[1;34m'
GREEN='\033[0;32m'
NC='\033[0m' # No Color

echo -e "${GREEN}Building webapi project...${NC}"
docker build -t webapigo .
echo -e "${BLUE}Open integration folder${NC}"
cd integration

echo -e "${GREEN}Build integration testing project...${NC}"
docker build -t webapitestgo .

echo -e "${BLUE}Go back main folder${NC}"
cd ..

echo -e "${GREEN}Integration testing is starting...${NC}"
docker-compose -f docker-compose.yml -f docker-compose.dev.yml up