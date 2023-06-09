**Prelaunch Configuration:**

1. Install Docker
2. Install weave, see `scripts/helper/weave_setup.sh`
3. Install Flannel, see `scripts/helper/flannel_setup.sh`
4. Configure `master_ip.txt`, just input the **IP address of master node**
5. Change mode of service, see `pkg/util/config/config.go`
6. Enjoy hacking!

**Usage:**

- Build k8s but not run
```
make build
```

- Build k8s and run as master and worker
```
make master
```

- Build k8s and run as worker
```
make worker
```

- Apply a Pod on one host
```
make testPod
```

- Apply a Pod, a Service on one host
```
make testService
```

- Apply a RS, a Service on one host
```
make testRS
```

- Apply a Stress Pod on one host
```
make testStress
```

- Apply a Pod, two Services, a gateway on one host
```
make testGateway
```

- Apply a Gpu Job on one host
```
make testGpu
```

- Apply a serverless test on one host
```
make testServerless
```

- Apply a microService test on one host
```
make testMicroService
```

- Stop k8s and clear k8s states
```
make clean
```

**CAUTION:**

- If you encounter any problem, try to run `make clean`, check your iptables if clean and run again.
The following command may help to clean iptables made by microService:
```
docker run --rm --name=istio-init --network=host --cap-add=NET_ADMIN istio/proxyv2:1.16.0 istio-clean-iptables
```
- Because our etcd is running in container, it has circle dependency with flannel, so it maybe fail when you run just after reboot. Just wait flannel online and rerun `make clean` and other command again.