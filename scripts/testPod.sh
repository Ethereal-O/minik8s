#!/bin/bash
./kubectl apply -f ./yaml/pod1.yaml
sleep 1
echo "[TestPod] Pod1 started!" 1>&2
