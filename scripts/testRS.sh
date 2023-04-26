#!/bin/bash
./kubectl apply -f ./yaml/rs1.yaml
sleep 1
echo "[TestRS] RS1 started!" 1>&2
./kubectl apply -f ./yaml/service1.yaml
sleep 1
echo "[TestRS] Service1 started!" 1>&2