#!/bin/bash

./kubectl apply -f ./yaml/GpuJobTest/gpu1.yaml
sleep 1
echo "[TestGpu] Gpu1(simple_change) started!" 1>&2
sleep 5
echo "[TestGpu] You may check the /home/shareDir to find that the target file exist!" 1>&2
sleep 20
echo "[TestGpu] You may check the /home/shareDir to check the result!(may take more time because of the server busy(may be more than one day))" 1>&2