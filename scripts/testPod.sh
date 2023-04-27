#!/bin/bash

# Start Pod1
./kubectl apply -f ./yaml/pod1.yaml
sleep 10
echo "[TestPod] Pod1 started!" 1>&2
./kubectl get -t Pod 1>&2

# Stop all Pause containers
docker stop $(docker ps -a | grep PAUSE | awk '{print $1}')
sleep 10
echo "[TestPod] Pod1 restarted!" 1>&2
./kubectl get -t Pod 1>&2

# Finally stop Pod1
./kubectl delete -t Pod -k Pod1
sleep 5
echo "[TestPod] Pod1 deleted!" 1>&2
./kubectl get -t Pod 1>&2
