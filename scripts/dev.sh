#!/usr/bin/env bash

set -e

GREEN='\033[0;32m'
NC='\033[0m'

echo -e "${GREEN}Starting dev env...${NC}"

docker compose up -d db

echo -e "${GREEN}Waiting for database...${NC}"
until pg_isready -h localhost -p 5432 -U myuser > /dev/null 2>&1; do
    sleep 1
done

echo -e "${GREEN}Running migrations...${NC}"
make -C backend migrate-up

echo -e "${GREEN}> Starting backend...${NC}"
cd backend
SESSION_KEY="dvamvtklzvqvfllotjkdrnnhtoqcdosa" DATABASE_URL="postgres://postgres:postgres@localhost:5432/tepia?sslmode=disable" air &

echo -e "${GREEN}> Starting frontend...${NC}"
cd ../frontend
# uncomment after npm
#npm install
#npm run dev &

wait
