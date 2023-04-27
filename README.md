**Prelaunch Configuration:**

1. Install Docker
2. Install weave, see `scripts/helper/weave_setup.sh`
3. Install ETCD, see `https://blog.csdn.net/qq_42874635/article/details/126906174`
4. Install NSQ, see `https://nsq.io/overview/quick_start.html`
5. Configure `master_ip.txt`, just input the **IP address of master node**
6. Enjoy hacking!

**Usage:**

- Build k8s but not run

```
make build
```

- Build k8s, run a master, a worker and a Pod on one host

```
make testPod
```

- Build k8s, run a master and a worker on one host
- Due to resource constraints (3 cloud servers), we have to run a master and a worker on one host

```
make master
```

- Build k8s, run a worker on one host

```
make worker
```

- Stop k8s and clear k8s states

```
make clean
```

- Stop everything, include k8s and environment

```
make stop
```