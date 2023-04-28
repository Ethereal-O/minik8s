package apiServer

import (
	"encoding/json"
	"fmt"
	"github.com/labstack/echo"
	ec "go.etcd.io/etcd/client/v3"
	"minik8s/pkg/etcd"
	"minik8s/pkg/messging"
	"minik8s/pkg/object"
	"minik8s/pkg/util/config"
	"minik8s/pkg/util/counter"
	"minik8s/pkg/util/stringParse"
	"net/http"
)

//--------------------- Basic Post Handler ---------------------------

func basic_post(c2 echo.Context) error {
	key := c2.FormValue("key")
	prix := c2.FormValue("prix")
	crt := c2.FormValue("crt")
	monitorKey := stringParse.Reform(key)
	prixFlag := false
	if prix == "true" {
		prixFlag = true
	}

	if crt == "" {
		newCrt := counter.GetMonitorCrt()
		monitorNum := monitorMap.Put(monitorKey, newCrt)
		if monitorNum == 1 {
			c := make(chan *ec.Event, 20)
			etcdStopFunc := etcd.Watch_etcd(key, prixFlag, c)
			producerStopFunc := messging.Producer(key, c)
			monitorEtcdStopMap[monitorKey] = etcdStopFunc
			monitorProducerStopMap[monitorKey] = producerStopFunc
		}
		return c2.String(http.StatusOK, newCrt)
	} else {
		monitorNum := monitorMap.Get(crt)
		if monitorNum == 0 {
			etcdStopFunc := monitorEtcdStopMap[monitorKey]
			producerStopFunc := monitorProducerStopMap[monitorKey]
			delete(monitorEtcdStopMap, monitorKey)
			delete(monitorProducerStopMap, monitorKey)
			etcdStopFunc()
			producerStopFunc()
		}
		return c2.String(http.StatusOK, "")
	}
}

//--------------------- Pod Handler ---------------------------

func pod_put(c echo.Context) error {
	podObject := new(object.Pod)
	if err := c.Bind(podObject); err != nil {
		return err
	}
	key := c.Request().RequestURI
	if podObject.Runtime.Uuid == "" {
		uuid := counter.GetUuid()
		podObject.Runtime.Uuid = uuid
	}
	if podObject.Runtime.Status == "" {
		podObject.Runtime.Status = config.CREATED_STATUS
	}
	if podObject.Runtime.ClusterIp == "" {
		podObject.Runtime.ClusterIp = NewPodIP()
	}

	pod, err := json.Marshal(podObject)
	if err != nil {
		fmt.Println(err.Error())
		return c.String(http.StatusInternalServerError, err.Error())
	}
	if err2 := etcd.Set_etcd(key, string(pod)); err2 != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return c.String(http.StatusOK, "ok")
}

func pod_get(c echo.Context) error {
	key := c.Request().RequestURI
	if c.Param("key") == config.EMPTY_FLAG {
		res := etcd.Get_etcd(key[0:len(key)-len(config.EMPTY_FLAG)], true)
		return c.JSON(http.StatusOK, res)
	} else {
		res := etcd.Get_etcd(key, false)
		return c.JSON(http.StatusOK, res)
	}
}

func pod_delete(c echo.Context) error {
	key := c.Request().RequestURI
	res := etcd.Get_etcd(key, false)
	if len(res) != 1 {
		return c.String(http.StatusInternalServerError, "not exist!")
	}
	var podObject object.Pod
	err := json.Unmarshal([]byte(res[0]), &podObject)
	if err != nil {
		return c.String(http.StatusInternalServerError, "unmarshal error!")
	}
	podObject.Runtime.Status = config.EXIT_STATUS
	pod, err := json.Marshal(podObject)
	if err != nil {
		fmt.Println(err.Error())
		return c.String(http.StatusInternalServerError, err.Error())
	}
	if err2 := etcd.Set_etcd(key, string(pod)); err2 != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.String(http.StatusOK, "delete successfully!")
}

//--------------------- ReplicaSet Handler ---------------------------

func replicaset_put(c echo.Context) error {
	rsObject := new(object.ReplicaSet)
	if err := c.Bind(rsObject); err != nil {
		return err
	}
	key := c.Request().RequestURI
	if rsObject.Runtime.Uuid == "" {
		uuid := counter.GetUuid()
		rsObject.Runtime.Uuid = uuid
	}
	if rsObject.Runtime.Status == "" {
		rsObject.Runtime.Status = config.RUNNING_STATUS
	}
	rs, err := json.Marshal(rsObject)
	if err != nil {
		fmt.Println(err.Error())
		return c.String(http.StatusInternalServerError, err.Error())
	}
	if err2 := etcd.Set_etcd(key, string(rs)); err2 != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return c.String(http.StatusOK, "ok")
}

func replicaset_get(c echo.Context) error {
	key := c.Request().RequestURI
	if c.Param("key") == config.EMPTY_FLAG {
		res := etcd.Get_etcd(key[0:len(key)-len(config.EMPTY_FLAG)], true)
		return c.JSON(http.StatusOK, res)
	} else {
		res := etcd.Get_etcd(key, false)
		return c.JSON(http.StatusOK, res)
	}
}

func replicaset_delete(c echo.Context) error {
	key := c.Request().RequestURI
	res := etcd.Get_etcd(key, false)
	if len(res) != 1 {
		return c.String(http.StatusInternalServerError, "not exist!")
	}
	var rsObject object.ReplicaSet
	err := json.Unmarshal([]byte(res[0]), &rsObject)
	if err != nil {
		return c.String(http.StatusInternalServerError, "unmarshal error!")
	}
	rsObject.Runtime.Status = config.EXIT_STATUS
	rs, err := json.Marshal(rsObject)
	if err != nil {
		fmt.Println(err.Error())
		return c.String(http.StatusInternalServerError, err.Error())
	}
	if err2 := etcd.Set_etcd(key, string(rs)); err2 != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	unbind(rsObject.Metadata.Name)
	return c.String(http.StatusOK, "delete successfully!")
}

//--------------------- AutoScaler Handler ---------------------------

func autoscaler_put(c echo.Context) error {
	hpaObject := new(object.AutoScaler)
	if err := c.Bind(hpaObject); err != nil {
		return err
	}
	key := c.Request().RequestURI
	if hpaObject.Runtime.Uuid == "" {
		uuid := counter.GetUuid()
		hpaObject.Runtime.Uuid = uuid
	}
	if hpaObject.Runtime.Status == "" {
		hpaObject.Runtime.Status = config.CREATED_STATUS
	}
	hpa, err := json.Marshal(hpaObject)
	if err != nil {
		fmt.Println(err.Error())
		return c.String(http.StatusInternalServerError, err.Error())
	}
	if err2 := etcd.Set_etcd(key, string(hpa)); err2 != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return c.String(http.StatusOK, "ok")
}

func autoscaler_get(c echo.Context) error {
	key := c.Request().RequestURI
	if c.Param("key") == config.EMPTY_FLAG {
		res := etcd.Get_etcd(key[0:len(key)-len(config.EMPTY_FLAG)], true)
		return c.JSON(http.StatusOK, res)
	} else {
		res := etcd.Get_etcd(key, false)
		return c.JSON(http.StatusOK, res)
	}
}

func autoscaler_delete(c echo.Context) error {
	key := c.Request().RequestURI
	res := etcd.Get_etcd(key, false)
	if len(res) != 1 {
		return c.String(http.StatusInternalServerError, "not exist!")
	}
	var hpaObject object.AutoScaler
	err := json.Unmarshal([]byte(res[0]), &hpaObject)
	if err != nil {
		return c.String(http.StatusInternalServerError, "unmarshal error!")
	}
	hpaObject.Runtime.Status = config.EXIT_STATUS
	hpa, err := json.Marshal(hpaObject)
	if err != nil {
		fmt.Println(err.Error())
		return c.String(http.StatusInternalServerError, err.Error())
	}
	if err2 := etcd.Set_etcd(key, string(hpa)); err2 != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return c.String(http.StatusOK, "delete successfully!")
}

//--------------------- Node Handler ---------------------------

func node_put(c echo.Context) error {
	nodeObject := new(object.Node)
	if err := c.Bind(nodeObject); err != nil {
		return err
	}
	key := c.Request().RequestURI
	if nodeObject.Runtime.Uuid == "" {
		uuid := counter.GetUuid()
		nodeObject.Runtime.Uuid = uuid
	}
	if nodeObject.Runtime.Status == "" {
		nodeObject.Runtime.Status = config.CREATED_STATUS
	}
	if nodeObject.Runtime.ClusterIp == "" {
		nodeObject.Runtime.ClusterIp = NewNodeIP()
	}
	node, err := json.Marshal(nodeObject)
	if err != nil {
		fmt.Println(err.Error())
		return c.String(http.StatusInternalServerError, err.Error())
	}
	if err2 := etcd.Set_etcd(key, string(node)); err2 != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return c.String(http.StatusOK, "ok")
}

func node_get(c echo.Context) error {
	key := c.Request().RequestURI
	if c.Param("key") == config.EMPTY_FLAG {
		res := etcd.Get_etcd(key[0:len(key)-len(config.EMPTY_FLAG)], true)
		return c.JSON(http.StatusOK, res)
	} else {
		res := etcd.Get_etcd(key, false)
		return c.JSON(http.StatusOK, res)
	}
}

func node_delete(c echo.Context) error {
	key := c.Request().RequestURI
	res := etcd.Get_etcd(key, false)
	if len(res) != 1 {
		return c.String(http.StatusInternalServerError, "not exist!")
	}
	var nodeObject object.Node
	err := json.Unmarshal([]byte(res[0]), &nodeObject)
	if err != nil {
		return c.String(http.StatusInternalServerError, "unmarshal error!")
	}
	nodeObject.Runtime.Status = config.EXIT_STATUS
	node, err := json.Marshal(nodeObject)
	if err != nil {
		fmt.Println(err.Error())
		return c.String(http.StatusInternalServerError, err.Error())
	}
	if err2 := etcd.Set_etcd(key, string(node)); err2 != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.String(http.StatusOK, "delete successfully!")
}

//--------------------- Service Handler ---------------------------

func service_put(c echo.Context) error {
	serviceObject := new(object.Service)
	if err := c.Bind(serviceObject); err != nil {
		return err
	}
	key := c.Request().RequestURI
	// create
	if serviceObject.Runtime.Uuid == "" {
		uuid := counter.GetUuid()
		serviceObject.Runtime.Uuid = uuid
	}
	if serviceObject.Runtime.Status == "" {
		serviceObject.Runtime.Status = config.RUNNING_STATUS
	}
	if serviceObject.Runtime.ClusterIp == "" {
		serviceObject.Runtime.ClusterIp = NewServiceIP()
	}
	service, err := json.Marshal(serviceObject)
	if err != nil {
		fmt.Println(err.Error())
		return c.String(http.StatusInternalServerError, err.Error())
	}
	if err2 := etcd.Set_etcd(key, string(service)); err2 != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return c.String(http.StatusOK, "ok")
}

func service_get(c echo.Context) error {
	key := c.Request().RequestURI
	if c.Param("key") == config.EMPTY_FLAG {
		res := etcd.Get_etcd(key[0:len(key)-len(config.EMPTY_FLAG)], true)
		return c.JSON(http.StatusOK, res)
	} else {
		res := etcd.Get_etcd(key, false)
		return c.JSON(http.StatusOK, res)
	}
}

func service_delete(c echo.Context) error {
	key := c.Request().RequestURI
	//err := etcd.Del_etcd(key)
	//if err != nil {
	//	return c.String(http.StatusInternalServerError, "delete failed!")
	//}
	res := etcd.Get_etcd(key, false)
	if len(res) != 1 {
		return c.String(http.StatusInternalServerError, "not exist!")
	}
	var serviceObject object.Service
	err := json.Unmarshal([]byte(res[0]), &serviceObject)
	if err != nil {
		return c.String(http.StatusInternalServerError, "unmarshal error!")
	}
	serviceObject.Runtime.Status = config.EXIT_STATUS
	service, err := json.Marshal(serviceObject)
	if err != nil {
		fmt.Println(err.Error())
		return c.String(http.StatusInternalServerError, err.Error())
	}
	if err2 := etcd.Set_etcd(key, string(service)); err2 != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.String(http.StatusOK, "delete successfully!")
}

//--------------------- Runtime Service Handler ---------------------------

func runtimeService_put(c echo.Context) error {
	runtimeServiceObject := new(object.RuntimeService)
	if err := c.Bind(runtimeServiceObject); err != nil {
		return err
	}
	key := c.Request().RequestURI
	// create
	if runtimeServiceObject.Service.Runtime.Uuid == "" {
		uuid := counter.GetUuid()
		runtimeServiceObject.Service.Runtime.Uuid = uuid
	}
	if runtimeServiceObject.Service.Runtime.Status == "" {
		runtimeServiceObject.Service.Runtime.Status = config.RUNNING_STATUS
	}
	serviceStatus, err := json.Marshal(runtimeServiceObject)
	if err != nil {
		fmt.Println(err.Error())
		return c.String(http.StatusInternalServerError, err.Error())
	}
	if err2 := etcd.Set_etcd(key, string(serviceStatus)); err2 != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return c.String(http.StatusOK, "ok")
}

func runtimeService_get(c echo.Context) error {
	key := c.Request().RequestURI
	fmt.Println(key)
	if c.Param("key") == config.EMPTY_FLAG {
		res := etcd.Get_etcd(key[0:len(key)-len(config.EMPTY_FLAG)], true)
		return c.JSON(http.StatusOK, res)
	} else {
		res := etcd.Get_etcd(key, false)
		return c.JSON(http.StatusOK, res)
	}
}

func runtimeService_delete(c echo.Context) error {
	key := c.Request().RequestURI
	res := etcd.Get_etcd(key, false)
	if len(res) != 1 {
		return c.String(http.StatusInternalServerError, "not exist!")
	}
	var runtimeServiceObject object.RuntimeService
	err := json.Unmarshal([]byte(res[0]), &runtimeServiceObject)
	if err != nil {
		return c.String(http.StatusInternalServerError, "unmarshal error!")
	}
	runtimeServiceObject.Service.Runtime.Status = config.EXIT_STATUS
	// because we unlock the lock in the runtimeService_put, so we don't need to unlock it here
	runtimeService, err := json.Marshal(runtimeServiceObject)
	if err != nil {
		fmt.Println(err.Error())
		return c.String(http.StatusInternalServerError, err.Error())
	}
	if err2 := etcd.Set_etcd(key, string(runtimeService)); err2 != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return c.String(http.StatusOK, "delete successfully!")
}
