#!/bin/bash

sh ./scripts/helper/weave_start.sh "subnet.txt" "worker"
if [ "$?" = 1 ]; then
  echo "[Worker] Failed to start weave subnet!" 1>&2
  exit 1
else
  echo "[Worker] Weave subnet started!" 1>&2
fi

# Get master IP
master_ip=""
while read -r line; do
  if echo "$line" | grep -q "master "; then
    master_ip=$(echo "$line" | awk '{print $2}')
    break
  fi
done < "subnet.txt"

if pgrep nsqd > /dev/null; then
  echo "[Worker] NSQ consumer is already running!" 1>&2
else
  nsqd --lookupd-tcp-address="$master_ip:4160" > nsqd.log 2>&1 &
  sleep 1
  echo "[Worker] NSQ consumer started!" 1>&2
fi

./kubectl worker > worker.log 2>&1 &
sleep 1
echo "[Worker] Worker node started!" 1>&2

echo "[Worker] Init finished!" 1>&2
