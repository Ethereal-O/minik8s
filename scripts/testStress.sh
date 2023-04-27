#!/bin/bash
./kubectl apply -f ./yaml/stress.yaml
sleep 1
echo "[TestStress] Stress started!" 1>&2