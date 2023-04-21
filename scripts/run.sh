etcd -listen-client-urls="http://0.0.0.0:2379" --advertise-client-urls="http://0.0.0.0:2379" > etcd.log 2>&1 &
sleep 2
./kubectl master > master.log 2>&1 &
sleep 2
./kubectl worker > worker.log 2>&1 &
sleep 2