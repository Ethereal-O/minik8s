package services

import (
	"encoding/json"
	"fmt"
	"minik8s/pkg/exeFile"
	"minik8s/pkg/object"
)

func (dnsManager *DnsManager) InitDnsTemplate() {
	dnsManager.initDnsReplicaSetTemplate()
	dnsManager.initDnsServiceTemplate()
	dnsManager.initGateWayReplicaSetTemplate()
	dnsManager.initGateWayServiceTemplate()
}

func (dnsManager *DnsManager) initDnsReplicaSetTemplate() {
	value, _, _ := exeFile.ReadYaml(DNS_REPLICATESET_TEMPLATE_FILEPATH)
	var ReplicaSetObject object.ReplicaSet
	err := json.Unmarshal([]byte(value), &ReplicaSetObject)
	if err != nil {
		fmt.Println("InitDnsReplicaSetTemplate fail" + err.Error())
		return
	}
	dnsManager.DnsTemplates.DnsReplicaSetTemplate = ReplicaSetObject
}

func (dnsManager *DnsManager) initDnsServiceTemplate() {
	value, _, _ := exeFile.ReadYaml(DNS_SERVICE_TEMPLATE_FILEPATH)
	var serviceObject object.Service
	err := json.Unmarshal([]byte(value), &serviceObject)
	if err != nil {
		fmt.Println("InitDnsServiceTemplate fail" + err.Error())
		return
	}
	dnsManager.DnsTemplates.DnsServiceTemplate = serviceObject
}

func (dnsManager *DnsManager) initGateWayReplicaSetTemplate() {
	value, _, _ := exeFile.ReadYaml(GATEWAY_REPLICATESET_TEMPLATE_FILEPATH)
	var replicaSetObject object.ReplicaSet
	err := json.Unmarshal([]byte(value), &replicaSetObject)
	if err != nil {
		fmt.Println("InitGateWayReplicaSetTemplate fail" + err.Error())
		return
	}
	dnsManager.DnsTemplates.GateWayReplicaSetTemplate = replicaSetObject
}

func (dnsManager *DnsManager) initGateWayServiceTemplate() {
	value, _, _ := exeFile.ReadYaml(GATEWAY_SERVICE_TEMPLATE_FILEPATH)
	var serviceObject object.Service
	err := json.Unmarshal([]byte(value), &serviceObject)
	if err != nil {
		fmt.Println("InitGateWayServiceTemplate fail" + err.Error())
		return
	}
	dnsManager.DnsTemplates.GateWayServiceTemplate = serviceObject
}

func GetGateWayReplicaSet(gatewayName string) object.ReplicaSet {
	template := dnsManager.DnsTemplates.GateWayReplicaSetTemplate
	template.Metadata.Name = GATEWAY_REPLICASET_PREFIX + gatewayName
	template.Spec.Template.Metadata.Labels[DNS_GATEWAY_SELECTOR] = gatewayName
	template.Spec.Template.Metadata.Name = GATEWAY_POD_PREFIX + gatewayName
	template.Spec.Template.Spec.Volumes[0].Path = NGINX_PATH_PREFIX + "/" + gatewayName
	template.Spec.Template.Spec.Containers[0].Name = GATEWAY_CONTAINER_PREFIX + gatewayName
	return template
}

func GetDnsService() object.Service {
	template := dnsManager.DnsTemplates.DnsServiceTemplate
	return template
}

func GetDnsReplicaSet() object.ReplicaSet {
	template := dnsManager.DnsTemplates.DnsReplicaSetTemplate
	return template
}

func GetGateWayService(gatewayName string) object.Service {
	template := dnsManager.DnsTemplates.GateWayServiceTemplate
	template.Metadata.Name = GATEWAY_SERVICE_PREFIX + gatewayName
	template.Spec.Selector[DNS_GATEWAY_SELECTOR] = gatewayName
	return template
}
