#!/bin/bash
# Should be used after testGateway finished!
sudo killall -9 -e ./kubectl master > /dev/null 2>&1
sleep 5
echo "[TestConsistency] Control plane killed!" 1>&2
./kubectl get -t Pod 1>&2
sudo ./kubectl master > master.log 2>&1 &
sleep 5
echo "[TestConsistency] Control plane restored!" 1>&2
./kubectl get -t Pod 1>&2
curl minik8s.com/service/ 1>&2