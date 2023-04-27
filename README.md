**Prelaunch Configuration:**

1. Install Docker
2. Install weave, see `scripts/helper/weave_setup.sh`
3. Install ETCD, see `https://blog.csdn.net/qq_42874635/article/details/126906174`
4. Install NSQ, see `https://nsq.io/overview/quick_start.html`
5. Install Prometheus, see `https://blog.csdn.net/qzcsu/article/details/124770699` (version:2.43.0)
6. Configure `master_ip.txt`, just input the **IP address of master node**
7. Enjoy hacking!

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

- Build k8s, run a master, a worker, a Pod on one host
```
make testPod
```

- Build k8s, run a master, a worker, a Pod, a Service on one host
```
make testService
```

- Build k8s, run a master, a worker, a RS, a Service on one host
```
make testRS
```

- Build k8s, run a master, a worker, a Stress Pod on one host
```
make testStress
```

- Stop k8s and clear k8s states
```
make clean
```