#!/bin/bash

# Start RS1
./kubectl apply -f ./yaml/rs1.yaml
sleep 80
echo "[TestRS] RS1 started!" 1>&2
./kubectl get -t Pod 1>&2

# Stop all downloader containers
docker stop $(docker ps -a | grep downloader | awk '{print $1}')
sleep 40
echo "[TestRS] RS1 restored its replicas!" 1>&2
./kubectl get -t Pod 1>&2