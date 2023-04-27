#!/bin/bash
./kubectl apply -f ./yaml/stress.yaml
sleep 20
echo "[TestStress] Stress started!" 1>&2
docker stats 1>&2