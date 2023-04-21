pids=$(pgrep -f "kubectl")

if [ -n "$pids" ]; then
    for pid in $pids; do
        kill "$pid"
    done
fi

pids=$(pgrep -f "etcd")

if [ -n "$pids" ]; then
    for pid in $pids; do
        kill "$pid"
    done
fi

rm -f ./kubectl
rm -rf ./default.etcd