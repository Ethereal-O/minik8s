**Usage:**

- Prepare environment for k8s, **only call it once after your VM startup!**
```
make prepare
```
- Build k8s but not run
```
make build
```
- Build k8s and run, but no kubectl commands
```
make run
```
- Build k8s and run testPod.sh
```
make testPod
```
- Stop k8s and clear its states (etcd, containers, ...)
```
make clean
```
- Stop everything, including k8s and environment
```
make stop
```