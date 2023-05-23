#!/bin/bash

# Start flannel network
grep -E 'SERVICE_POLICY *= SERVICE_POLICY_IPTABLES|SERVICE_POLICY *= SERVICE_POLICY_MICROSERVICE' ./pkg/util/config/config.go > /dev/null 2>&1
if [ "$?" != 1 ]; then
  bash ./scripts/helper/flannel_start.sh
  if [ "$?" = 1 ]; then
    echo "[Worker] Failed to start flannel network!" 1>&2
    exit 1
  else
    echo "[Worker] Flannel network started!" 1>&2
  fi
fi