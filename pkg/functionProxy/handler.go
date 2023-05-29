package functionProxy

import (
	"encoding/json"
	"fmt"
	"github.com/labstack/echo"
	"minik8s/pkg/client"
	"minik8s/pkg/object"
	"minik8s/pkg/util/config"
	"minik8s/pkg/util/tools"
	"net/http"
	"strconv"
	"strings"
)

//----------------------- Handler ---------------------

func forwardRequest(c2 echo.Context) error {
	// Step1: Parse the formdata of the request
	formParams, _ := c2.FormParams()
	formData := make(map[string]string)
	for key, values := range formParams {
		value := strings.Join(values, ",")
		formData[key] = value
	}

	// Step2: Check the state of the function
	funcName := formData["function"]
	tarFunction := client.GetFunction(funcName)
	// If the function doesn't exist, just do nothing and return the corresponding information
	if tarFunction == nil || tarFunction.Runtime.Status == config.EXIT_STATUS {
		return c2.String(http.StatusOK, "Function not exist!")
	}
	// If the function exist but scale-to-0, activate it and return the corresponding information
	if tarFunction.Runtime.Status != config.RUNNING_STATUS {
		activate(config.FUNC_NAME + "_rs_" + tarFunction.FaasName)
		return c2.String(http.StatusOK, "Function not up!")
	}

	// Step3: Request to the real server or from cache and return the result
	addFlow(config.FUNC_NAME + "_service_" + tarFunction.FaasName)
	formData["module"] = tarFunction.Module
	if cache.GetCache(tools.MD5(formData)) != "" {
		fmt.Println("From cache")
		res := cache.GetCache(tools.MD5(formData))
		return c2.String(http.StatusOK, res)
	} else {
		fmt.Println("From server")
		url := "http://" + tarFunction.Runtime.FunctionIp + ":8081" + "/run"
		res := client.ForwardPostData(formData, url)
		cache.PutCache(tools.MD5(formData), res)
		return c2.String(http.StatusOK, res)
	}
}

func doWorkflow(c2 echo.Context) error {
	// Step1: Parse the formdata of the request
	// There are only two pair of kv:the workflow and the initial params
	// Get them and unmarshal into the corresponding data structure
	formParams, _ := c2.FormParams()
	formData := make(map[string]string)
	for key, values := range formParams {
		value := strings.Join(values, ",")
		formData[key] = value
	}
	workflow_json := formData["workflow"]
	param_json := formData["params"]
	workflow := object.WorkFlow{}
	param := make(map[string]string)
	_ = json.Unmarshal([]byte(workflow_json), &workflow)
	_ = json.Unmarshal([]byte(param_json), &param)

	// Step2: Check the state of the function
	// If any function in dag doesn't exist, just do nothing and return the corresponding information
	for _, node := range workflow.Nodes {
		tarFunction := client.GetFunction(node.FuncName)
		if tarFunction == nil {
			return c2.String(http.StatusOK, "Function not exist!")
		}
	}
	// Active all function in dag
	coldStartFlag := false
	for _, node := range workflow.Nodes {
		tarFunction := client.GetFunction(node.FuncName)
		if tarFunction.Runtime.Status != config.RUNNING_STATUS {
			activate(config.FUNC_NAME + "_rs_" + tarFunction.FaasName)
			coldStartFlag = true
		}
	}
	// If any function in dag exist but scale-to-0, activate it and return the corresponding information
	if coldStartFlag {
		return c2.String(http.StatusOK, "Function not up!")
	}

	// Step3: Request to the real server or from cache recursively and return the result
	tmpNode := getNode(workflow.StartNode, workflow)
	tmpParam := param
	triggerPath := workflow.StartNode
	var res string
	for tmpNode.NodeName != "" {
		// Record the trigger path
		if tmpNode.NodeName != workflow.StartNode {
			triggerPath += "-" + tmpNode.NodeName
		}
		// Do the request
		tarFunction := client.GetFunction(tmpNode.FuncName)
		addFlow(config.FUNC_NAME + "_service_" + tarFunction.FaasName)
		tmpParam["module"] = tarFunction.Module
		tmpParam["function"] = tarFunction.FuncName
		if cache.GetCache(tools.MD5(tmpParam)) != "" {
			fmt.Println("From cache")
			res = cache.GetCache(tools.MD5(tmpParam))
		} else {
			fmt.Println("From server")
			url := "http://" + tarFunction.Runtime.FunctionIp + ":8081" + "/run"
			res = client.ForwardPostData(tmpParam, url)
			cache.PutCache(tools.MD5(tmpParam), res)
		}
		//Prepare the param and choose the next dag node
		_ = json.Unmarshal([]byte(res), &tmpParam)
		tmpNode = selectNode(tmpNode, tmpParam, workflow)
	}

	// Step4: Complete the ans and return to the client
	resMap := make(map[string]string)
	_ = json.Unmarshal([]byte(res), &resMap)
	delete(resMap, "function")
	delete(resMap, "module")
	resMap["trigger_path"] = triggerPath
	resJson, _ := json.Marshal(resMap)
	return c2.String(http.StatusOK, string(resJson))
}

// ----------------------- Function ---------------------
func getNode(nodeName string, workflow object.WorkFlow) object.DagNode {
	for _, node := range workflow.Nodes {
		if node.NodeName == nodeName {
			return node
		}
	}
	return object.DagNode{}
}

func selectNode(node object.DagNode, param map[string]string, workflow object.WorkFlow) object.DagNode {
	if len(node.Choices) == 0 {
		return object.DagNode{}
	}

	for _, choice := range node.Choices {
		tarVariable := choice.Condition.TarVariable
		tarValue := choice.Condition.TarValue
		relation := choice.Condition.Relation
		nextNode := choice.NextNode
		if choice.Condition.Relation == "" {
			return getNode(nextNode, workflow)
		}
		if _, ok := param[tarVariable]; ok {
			variable, _ := strconv.Atoi(param[tarVariable])
			value, _ := strconv.Atoi(tarValue)
			if relation == "eq" && variable == value {
				return getNode(nextNode, workflow)
			} else if relation == "ne" && variable != value {
				return getNode(nextNode, workflow)
			} else if relation == "le" && variable <= value {
				return getNode(nextNode, workflow)
			} else if relation == "lt" && variable < value {
				return getNode(nextNode, workflow)
			} else if relation == "ge" && variable >= value {
				return getNode(nextNode, workflow)
			} else if relation == "gt" && variable > value {
				return getNode(nextNode, workflow)
			}
		}
	}
	return getNode(node.Choices[0].NextNode, workflow)
}
