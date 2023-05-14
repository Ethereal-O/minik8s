#!/bin/bash
sudo killall kubectl > /dev/null 2>&1
echo "[Cleaner] K8S Master/Worker stopped!" 1>&2
sudo killall etcd > /dev/null 2>&1
echo "[Cleaner] ETCD stopped!" 1>&2
sudo killall nsqlookupd > /dev/null 2>&1
sudo killall nsqd > /dev/null 2>&1
sudo killall nsqadmin > /dev/null 2>&1
echo "[Cleaner] NSQ Producer/Consumer stopped!" 1>&2
sudo weave reset > /dev/null 2>&1
echo "[Cleaner] Weave subnet stopped!" 1>&2

sudo rm -rf ./default.etcd
echo "[Cleaner] ETCD data cleared!" 1>&2
sudo rm -f ./*.dat
echo "[Cleaner] NSQ data cleared!" 1>&2
#sudo rm -rf /home/os/minik8s/DNS > /dev/null 2>&1
sudo rm -rf /home/os/minik8s/Gateway > /dev/null 2>&1
echo "[Cleaner] DNS data cleared!" 1>&2
sudo docker kill $(docker ps -q) > /dev/null 2>&1
sudo docker rm $(docker ps -aq) > /dev/null 2>&1
echo "[Cleaner] Containers cleared!" 1>&2

echo "[Cleaner] All states cleared!" 1>&2

