#!/bin/bash

sh ./scripts/helper/weave_start.sh
if [ "$?" = 1 ]; then
  echo "[Master] Failed to start weave subnet!" 1>&2
  exit 1
else
  echo "[Master] Weave subnet started!" 1>&2
fi

etcd -listen-client-urls="http://0.0.0.0:2379" --advertise-client-urls="http://0.0.0.0:2379" > etcd.log 2>&1 &
sleep 1
echo "[Master] ETCD started!" 1>&2

if pgrep nsqlookupd > /dev/null; then
  echo "[Master] NSQ producer is already running!" 1>&2
else
  nsqlookupd > nsqlookupd.log 2>&1 &
  sleep 1
  nsqadmin --lookupd-http-address=127.0.0.1:4161 > nsqadmin.log 2>&1 &
  sleep 1
  nsqd --lookupd-tcp-address=127.0.0.1:4160 > nsqd.log 2>&1 &
  sleep 1
  
  pubUrl="http://127.0.0.1:4151/pub?topic="
  curl -d "hello" "${pubUrl}Pod"
  curl -d "hello" "${pubUrl}ReplicaSet"
  curl -d "hello" "${pubUrl}AutoScaler"
  curl -d "hello" "${pubUrl}Service"
  curl -d "hello" "${pubUrl}RuntimeService"
  curl -d "hello" "${pubUrl}Node"
  curl -d "hello" "${pubUrl}DNS"
  curl -d "hello" "${pubUrl}Gateway"
  curl -d "hello" "${pubUrl}RuntimeGateway"
  curl -d "hello" "${pubUrl}GpuJob"
  curl -d "hello" "${pubUrl}ServerlessFunctions"
  curl -d "hello" "${pubUrl}Function"
  curl -d "hello" "${pubUrl}TransFile"

  echo "[Master] NSQ producer+consumer started!" 1>&2
fi

# here to copy file to /home/os/minik8s/DNS, so that we can use it in minik8s
sudo rm -rf /home/os/minik8s/DNS > /dev/null 2>&1
sudo rm -rf /home/os/minik8s/Gateway > /dev/null 2>&1
sudo mkdir -p /home/os/minik8s/DNS > /dev/null 2>&1
sudo mkdir -p /home/os/minik8s/Gateway > /dev/null 2>&1
sudo cp ./template/config/CORE_DNS_CONFIG/* /home/os/minik8s/DNS > /dev/null 2>&1
echo "[Master] DNS config created!" 1>&2

sudo ./kubectl master > master.log 2>&1 &
sleep 5
echo "[Master] Control plane started!" 1>&2

echo "[Master] Init finished!" 1>&2
