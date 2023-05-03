#!/bin/bash

# Start Pod1
./kubectl apply -f ./yaml/pod1.yaml
sleep 60
echo "[TestPod] Pod1 started!" 1>&2
./kubectl get -t Pod 1>&2

# Stop all downloader containers
docker stop $(docker ps -a | grep downloader | awk '{print $1}')
sleep 30
echo "[TestPod] Pod1 restarted!" 1>&2
./kubectl get -t Pod 1>&2

# Finally stop Pod1
./kubectl delete -t Pod -k Pod1
sleep 10
echo "[TestPod] Pod1 deleted!" 1>&2
./kubectl get -t Pod 1>&2
