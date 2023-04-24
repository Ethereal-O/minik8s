#!/bin/bash
./kubectl apply -f ./yaml/node1.yaml
sleep 1
echo "[TestPod] Node1 started!" 1>&2
./kubectl apply -f ./yaml/pod1.yaml
sleep 1
echo "[TestPod] Pod1 started!" 1>&2