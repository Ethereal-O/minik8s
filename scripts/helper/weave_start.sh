#!/bin/bash

# Read config file
config_file="$1"
if [ ! -f "$config_file" ]; then
  echo "Config file $config_file does not exist!" 1>&2
  exit 1
fi

# Get host IP and role
ip_addr=$(/sbin/ifconfig -a|grep inet|grep -v 127.0.0.1|grep -v inet6|awk '{print $2}'|tr -d "addr:")
role=""
subnet_ip=""
while read -r line; do
  if echo "$line" | grep -q "$ip_addr"; then
    role=$(echo "$line" | awk '{print $1}')
    subnet_ip=$(echo "$line" | awk '{print $3}')
    break
  fi
done < "$config_file"
if [ -z "$role" ]; then
  echo "Config file has no IP address: $ip_addr!" 1>&2
  exit 1
fi

# Get master IP
master_ip=""
while read -r line; do
  if echo "$line" | grep -q "master "; then
    master_ip=$(echo "$line" | awk '{print $2}')
    break
  fi
done < "$config_file"
if [ -z "$master_ip" ]; then
  echo "Config file has no master!" 1>&2
  exit 1
fi

# Start weave subnet based on role
if [ "$role" = "master" ] && [ "$2" = "master" ]; then
  weave reset > /dev/null 2>&1
  weave launch > /dev/null 2>&1
  weave expose "$subnet_ip/16" > /dev/null 2>&1
elif [ "$role" = "worker" ] && [ "$2" = "worker" ]; then
  weave reset > /dev/null 2>&1
  weave launch "$master_ip" > /dev/null 2>&1
  weave expose "$subnet_ip/16" > /dev/null 2>&1
fi


exit 0