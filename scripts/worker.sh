#!/bin/bash

# here to copy file to /home/os/minik8s/DNS, so that we can use it in minik8s
sudo rm -rf /home/os/minik8s/DNS > /dev/null 2>&1
sudo rm -rf /home/os/minik8s/Gateway > /dev/null 2>&1
sudo rm -rf /home/os/minik8s/Service > /dev/null 2>&1
sudo rm -rf /home/os/minik8s/Forward > /dev/null 2>&1
sudo mkdir -p /home/os/minik8s/DNS > /dev/null 2>&1
sudo mkdir -p /home/os/minik8s/Gateway > /dev/null 2>&1
sudo mkdir -p /home/os/minik8s/Service > /dev/null 2>&1
sudo mkdir -p /home/os/minik8s/Forward > /dev/null 2>&1
sudo cp -r ./template/config/CORE_DNS_CONFIG/* /home/os/minik8s/DNS > /dev/null 2>&1
sudo cp -r ./template/config/NGINX_TEMPLATE/* /home/os/minik8s/Forward > /dev/null 2>&1
sudo cp /etc/hosts /etc/hosts.bak > /dev/null 2>&1
echo "[Worker] DNS config created!" 1>&2

# Start weave subnet
sh ./scripts/helper/weave_start.sh
if [ "$?" = 1 ]; then
  echo "[Worker] Failed to start weave subnet!" 1>&2
  exit 1
else
  echo "[Worker] Weave subnet started!" 1>&2
fi

# Start nsq-consumer (nsqd)
host_ip=$(/sbin/ifconfig -a|grep 192.168|awk '{print $2}'|tr -d "addr:")
master_ip=$(cat master_ip.txt)
if [ "$host_ip" != "$master_ip" ]; then
  docker run -d \
    --network=host \
    --name nsq-consumer \
    -v "$(pwd)/scripts/helper/nsq_consumer.sh:/nsq_consumer.sh" \
    -v "$(pwd)/master_ip.txt:/master_ip.txt" \
    nsqio/nsq \
    sh /nsq_consumer.sh \
    > /dev/null 2>&1
  echo "[Worker] NSQ consumer started!" 1>&2
fi

# Start worker node
sudo ./kubectl worker > worker.log 2>&1 &
sleep 5
echo "[Worker] Worker node started!" 1>&2

echo "[Worker] Init finished!" 1>&2
