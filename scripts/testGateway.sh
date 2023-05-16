#!/bin/bash
sleep 5
./kubectl apply -f ./yaml/pod1.yaml
sleep 1
echo "[TestGateway] Pod1 started!" 1>&2
./kubectl apply -f ./yaml/pod2.yaml
sleep 1
echo "[TestGateway] Pod2 started!" 1>&2
./kubectl apply -f ./yaml/service1.yaml
sleep 1
echo "[TestGateway] Service1 started!" 1>&2
./kubectl apply -f ./yaml/service2.yaml
sleep 1
echo "[TestGateway] Service2 started!" 1>&2
./kubectl apply -f ./yaml/gateway1.yaml
sleep 1
echo "[TestGateway] Gateway1 started!" 1>&2