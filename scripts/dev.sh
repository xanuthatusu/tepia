#!/usr/bin/env bash

set -e

GREEN='\033[0;32m'
NC='\033[0m'

echo -e "${GREEN}Starting dev env...${NC}"

echo -e "${GREEN}> Starting backend...${NC}"
cd backend
air &

echo -e "${GREEN}> Starting frontend...${NC}"
cd ../frontend
# uncomment after npm
#npm install
#npm run dev &

wait
