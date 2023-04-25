#!/bin/bash
sudo killall kubectl > /dev/null 2>&1
echo "[Cleaner] K8S Master/Worker stopped!" 1>&2
sudo killall etcd > /dev/null 2>&1
echo "[Cleaner] ETCD stopped!" 1>&2
sudo weave reset > /dev/null 2>&1
echo "[Cleaner] Weave subnet stopped!" 1>&2

sudo rm -rf ./default.etcd
echo "[Cleaner] ETCD data cleared!" 1>&2
sudo rm -f ./*.dat
echo "[Cleaner] NSQ data cleared!" 1>&2
sh ./scripts/helper/delete_all_containers.sh
echo "[Cleaner] Containers cleared!" 1>&2

echo "[Cleaner] All states cleared!" 1>&2