#!/bin/bash
set -e

echo "building docker images..."

docker build --build-arg CMD_PATH=setup -t localhost:32000/setup:latest .
docker build --build-arg CMD_PATH=txrmngr -t localhost:32000/txrmngr:latest .
docker build --build-arg CMD_PATH=operator -t localhost:32000/operator:latest .
docker build --build-arg CMD_PATH=watchdog -t localhost:32000/watchdog:latest .
docker build --build-arg CMD_PATH=service -t localhost:32000/service:latest .

echo "pushing docker images to local registry..."

docker push localhost:32000/setup:latest
docker push localhost:32000/txrmngr:latest
docker push localhost:32000/operator:latest
docker push localhost:32000/watchdog:latest
docker push localhost:32000/service:latest

echo "build and push completed."