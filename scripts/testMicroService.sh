#!/bin/bash
## IMPORTANT: This script should be run NOT as root
sleep 20
# Start Pod1 and Pod2
./kubectl apply -f ./yaml/MicroServiceTest/pod1.yaml
./kubectl apply -f ./yaml/MicroServiceTest/pod2.yaml
sleep 10
echo "[TestMicroService] Pod1 and Pod2 started!" 1>&2

# Start Service1
./kubectl apply -f ./yaml/MicroServiceTest/service1.yaml
sleep 20
echo "[TestMicroService] Service1 started!" 1>&2

# Start VirtualService1
./kubectl apply -f ./yaml/MicroServiceTest/virtualService1.yaml
sleep 5
echo "[TestMicroService] VirtualService1 started!" 1>&2

# test curl once
echo "[TestMicroService] Testing curl once..." 1>&2
curl 100.100.20.1

# test curl 100 times
echo "[TestMicroService] Testing curl 100 times..." 1>&2
for i in {1..100}
do
    curl 100.100.20.1 > /dev/null 2>&1
done

grep 100.100.20.1 ./worker.log | sort -rn | uniq -c | sort -rn
