**Subnet Configuration:**

Example of `subnet.txt`:

```
master 192.168.29.128 10.10.0.1
worker 192.168.29.129 10.10.0.2
worker 192.168.29.130 10.10.0.3
```

Each row is a host.
The first column refers to the role of the host, i.e., master or worker.
The second column refers to the public IP of the host, which should be seen in `ifconfig`.
The third column refers to the subnet IP of the host, which is valid only in the cluster.

For example, if you test minik8s on a single host, and the public IP of the host is `192.168.1.1`, then you should configure `subnet.txt`:

```
master 192.168.1.1 10.10.0.1
```

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