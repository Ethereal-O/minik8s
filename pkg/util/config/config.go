package config

const BASE_URL string = "192.168.29.132"
const APISERVER_URL string = "http://" + BASE_URL + ":8080"
const ETCD_Endpoints string = BASE_URL + ":2379"
const NSQ_PEODUCER string = BASE_URL + ":4150"
const NSQ_CONSUMER string = BASE_URL + ":4161"

var POD_TYPE = "Pod"
var REPLICASET_TYPE = "Replicaset"
var TP = []string{POD_TYPE, REPLICASET_TYPE}

var EMPTY_FLAG = "none"

var RUNNING_STATUS = "RUNNING"
var EXIT_STATUS = "EXIT"
