#!/bin/bash

# Start etcd
docker run -d \
  --env ALLOW_NONE_AUTHENTICATION=yes \
  --env ETCD_ENABLE_V2=true \
  --env ETCDCTL_API=2 \
  --network=host \
  --name etcd \
  bitnami/etcd \
  > /dev/null 2>&1
sleep 1
docker exec etcd \
  etcdctl --endpoints http://127.0.0.1:2379 set /coreos.com/network/config '{"Network": "172.17.0.0/16", "SubnetLen": 24, "SubnetMin": "172.17.1.0","SubnetMax": "172.17.20.0", "Backend": {"Type": "vxlan"}}' \
  > /dev/null 2>&1
echo "[Master] ETCD started!" 1>&2

# Start prometheus
docker run -d \
  --network=host \
  --name prometheus \
  -v "$(pwd)/prometheus.yml:/etc/prometheus/prometheus.yml" \
  prom/prometheus \
  > /dev/null 2>&1
echo "[Master] Prometheus started!" 1>&2

# Start nsq-producer (nsqlookupd + nsqadmin + nsqd)
docker run -d \
  --network=host \
  --name nsq-producer \
  -v "$(pwd)/scripts/helper/nsq_producer.sh:/nsq_producer.sh" \
  nsqio/nsq \
  sh /nsq_producer.sh \
  > /dev/null 2>&1
echo "[Master] NSQ producer started!" 1>&2

# Wake up topics
sleep 3
pubUrl="http://127.0.0.1:4151/pub?topic="
curl -d "hello" "${pubUrl}Pod" > /dev/null 2>&1
curl -d "hello" "${pubUrl}ReplicaSet" > /dev/null 2>&1
curl -d "hello" "${pubUrl}DaemonSet" > /dev/null 2>&1
curl -d "hello" "${pubUrl}AutoScaler" > /dev/null 2>&1
curl -d "hello" "${pubUrl}Service" > /dev/null 2>&1
curl -d "hello" "${pubUrl}RuntimeService" > /dev/null 2>&1
curl -d "hello" "${pubUrl}VirtualService" > /dev/null 2>&1
curl -d "hello" "${pubUrl}Node" > /dev/null 2>&1
curl -d "hello" "${pubUrl}DNS" > /dev/null 2>&1
curl -d "hello" "${pubUrl}Gateway" > /dev/null 2>&1
curl -d "hello" "${pubUrl}RuntimeGateway" > /dev/null 2>&1
curl -d "hello" "${pubUrl}GpuJob" > /dev/null 2>&1
curl -d "hello" "${pubUrl}ServerlessFunctions" > /dev/null 2>&1
curl -d "hello" "${pubUrl}Function" > /dev/null 2>&1
curl -d "hello" "${pubUrl}TransFile" > /dev/null 2>&1

# Start control plane
sudo ./kubectl master > master.log 2>&1 &
sleep 5
echo "[Master] Control plane started!" 1>&2

echo "[Master] Init finished!" 1>&2
