#!/bin/bash
./kubectl apply -f ./yaml/pod1.yaml
sleep 1
echo "[TestService] Pod1 started!" 1>&2
./kubectl apply -f ./yaml/service1.yaml
sleep 1
echo "[TestService] Service1 started!" 1>&2