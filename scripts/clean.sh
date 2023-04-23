killall etcd > /dev/null 2>&1
killall kubectl > /dev/null 2>&1

rm -rf ./default.etcd
sh ./scripts/helper/delete_all_containers.sh