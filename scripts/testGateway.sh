#!/bin/bash
./kubectl apply -f ./yaml/pod1.yaml
sleep 1
echo "[TestGateway] Pod1 started!" 1>&2
./kubectl apply -f ./yaml/service1.yaml
sleep 1
echo "[TestGateway] Service1 started!" 1>&2
./kubectl apply -f ./yaml/gateway1.yaml
sleep 1
echo "[TestGateway] Gateway1 started!" 1>&2