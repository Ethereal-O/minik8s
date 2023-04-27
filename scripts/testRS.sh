#!/bin/bash

# Start RS1
./kubectl apply -f ./yaml/rs1.yaml
sleep 30
echo "[TestRS] RS1 started!" 1>&2
./kubectl get -t Pod 1>&2

# Stop all Pause containers
docker stop $(docker ps -a | grep PAUSE | awk '{print $1}')
sleep 30
echo "[TestRS] RS1 restored its replicas!" 1>&2
./kubectl get -t Pod 1>&2