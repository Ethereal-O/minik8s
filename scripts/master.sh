#!/bin/bash

etcd -listen-client-urls="http://0.0.0.0:2379" --advertise-client-urls="http://0.0.0.0:2379" --enable-v2 > etcd.log 2>&1 &
sleep 1
ETCDCTL_API=2 etcdctl --endpoints http://127.0.0.1:2379 set /coreos.com/network/config '{"Network": "172.17.0.0/16", "SubnetLen": 24, "SubnetMin": "172.17.1.0","SubnetMax": "172.17.20.0", "Backend": {"Type": "vxlan"}}' > /dev/null 2>&1
echo "[Master] ETCD started!" 1>&2

prometheus --config.file=prometheus.yml > prometheus.log 2>&1 &
echo "[Master] Prometheus started!" 1>&2

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
  curl -d "hello" "${pubUrl}DaemonSet"
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

sudo ./kubectl master > master.log 2>&1 &
sleep 5
echo "[Master] Control plane started!" 1>&2

echo "[Master] Init finished!" 1>&2
