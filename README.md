# Prepare environment for running, only call it once after your VM startup!
make prepare

# Build k8s but not run
make build

# Build k8s and run, but no kubectl commands
make run

# Build k8s and run testPod.sh
make testPod

# Stop k8s and clear etcd and containers
make clean

# Stop everything, include k8s and environment
make stop