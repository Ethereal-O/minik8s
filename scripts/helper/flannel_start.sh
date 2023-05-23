#!/bin/bash

# Read config file
config_file="master_ip.txt"
if [ ! -f "$config_file" ]; then
  echo "Config file $config_file does not exist!" 1>&2
  exit 1
fi

# Get master IP
master_ip=$(cat "$config_file")
if [ -z "$master_ip" ]; then
  echo "Config file $config_file is empty!" 1>&2
  exit 1
fi

# Start flannel
sleep 1
flannel --etcd-endpoints="http://$master_ip:2379" --ip-masq=true > flannel.log 2>&1 &
sleep 4
sudo systemctl stop docker.socket > /dev/null 2>&1
sudo systemctl stop docker > /dev/null 2>&1
sudo killall dockerd > /dev/null 2>&1
export FLANNEL_SUBNET=172.17.0.1/16
export FLANNEL_MTU=1450
source /run/flannel/subnet.env
sudo echo -e "{\n\t\"bip\":\"${FLANNEL_SUBNET}\",\n\t\"mtu\":${FLANNEL_MTU}\n}" > /etc/docker/daemon.json
sudo systemctl start docker > /dev/null 2>&1
sleep 1

exit 0