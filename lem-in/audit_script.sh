#!/bin/bash

GREEN='\033[0;32m'
YELLOW='\033[0;33m'
NC='\033[0m' # No Color

for i in 0 1; do
    echo -e "$GREEN------ Running badexample0$i.txt ------$NC"
    echo
    go run . badexample0$i.txt
    echo -en "\n${YELLOW}Press enter to continue...${NC}"
    read
    echo
done

for i in {0..7}; do
    echo -e "$GREEN------ Running example0$i.txt ------$NC"
    echo
    go run . example0$i.txt
    echo -en "\n${YELLOW}Press enter to continue...${NC}"
    read
    echo
done