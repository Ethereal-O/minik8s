package config

const APISERVER_URL string = "http://192.168.142.135:8080"
const ETCD_Endpoints string = "192.168.142.135:2379"
const NSQ_PEODUCER string = "192.168.142.135:4150"
const NSQ_CONSUMER string = "192.168.142.135:4161"

var POD_TYPE = "Pod"
var REPLICASET_TYPE = "Replicaset"
var TP = []string{POD_TYPE, REPLICASET_TYPE}

var EMPTY_FLAG = "none"

var RUNNING_STATUS = "RUNNING"
var EXIT_STATUS = "EXIT"
