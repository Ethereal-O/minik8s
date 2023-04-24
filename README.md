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