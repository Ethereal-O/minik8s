#!/bin/bash
# Start RS1 and Service1
./kubectl apply -f ./yaml/rs1.yaml
./kubectl apply -f ./yaml/service1.yaml
sleep 80
echo "[TestService] RS1 started!" 1>&2
./kubectl get -t Pod 1>&2
./kubectl get -t Service 1>&2

# Stop all downloader containers
docker stop $(docker ps -a | grep downloader | awk '{print $1}')
sleep 40
echo "[TestService] RS1 restored its replicas!" 1>&2
./kubectl get -t Pod 1>&2
./kubectl get -t Service 1>&2

# Finally stop RS1
./kubectl delete -t ReplicaSet -k RS1
sleep 10
echo "[TestService] RS1 deleted!" 1>&2
./kubectl get -t Pod 1>&2
./kubectl get -t Service 1>&2