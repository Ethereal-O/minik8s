#!/bin/bash
./kubectl apply -f ./yaml/pod1.yaml
sleep 1
echo "[TestPod] Pod1 started!" 1>&2
./kubectl apply -f ./yaml/service.yaml
sleep 1
echo "[TestService] testService started!" 1>&2