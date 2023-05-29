#!/bin/bash

sleep 5
echo "[TestServerless] Register the functions!" 1>&2
./kubectl apply -f ./yaml/ServerlessTest/func1.yaml 1>&2
sleep 1

echo "[TestServerless] Request for single function!(cold start)" 1>&2
./kubectl request -f minus -p ./yaml/ServerlessTest/param1.yaml 1>&2
sleep 1

echo "[TestServerless] Request for single function!" 1>&2
./kubectl request -f minus -p ./yaml/ServerlessTest/param1.yaml 1>&2
./kubectl request -f add -p ./yaml/ServerlessTest/param1.yaml 1>&2
sleep 1

for i in {1..15}; do
  ./kubectl request -f add -p ./yaml/ServerlessTest/param1.yaml > /dev/null 2>&1
done
echo "[TestServerless] After a lot requests,the replicas should grow up!(wait 30s and see the result)" 1>&2
sleep 30
./kubectl get -t ReplicaSet 1>&2

echo "[TestServerless] Request for DAG workflow!" 1>&2
./kubectl request -d ./yaml/ServerlessTest/dag1.yaml -p ./yaml/ServerlessTest/param1.yaml 1>&2

echo "[TestServerless] After no requests,the replicas should scale-to-0!(wait 180s and see the result)" 1>&2
sleep 180
./kubectl get -t ReplicaSet 1>&2
