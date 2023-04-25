#!/bin/bash

# Read config file
config_file="master_ip.txt"
if [ ! -f "$config_file" ]; then
  echo "Config file $config_file does not exist!" 1>&2
  exit 1
fi

# Get host IP
host_ip=$(/sbin/ifconfig -a|grep 192.168|awk '{print $2}'|tr -d "addr:")

# Get master IP
master_ip=$(cat "$config_file")
if [ -z "$master_ip" ]; then
  echo "Config file $config_file is empty!" 1>&2
  exit 1
fi

# Start weave subnet based on role
if [ "$host_ip" = "$master_ip" ]; then
  weave launch > /dev/null 2>&1
else
  weave launch "$master_ip" > /dev/null 2>&1
fi

exit 0