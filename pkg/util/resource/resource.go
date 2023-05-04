package resource

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

func ConvertMemoryToBytes(memoryStr string) int64 {
	var bytes int64 = 1024 * 1024 * 200
	memoryStr = strings.TrimSpace(memoryStr)
	if memoryStr == "" {
		return bytes //default 200 MB
	}
	memoryStr = strings.ToLower(memoryStr)
	if memoryStr[len(memoryStr)-2:] == "gi" {
		val, err := strconv.ParseFloat(memoryStr[:len(memoryStr)-2], 64)
		if err != nil {
			fmt.Println(err.Error())
			return bytes
		}
		bytes = int64(math.Round(val * math.Pow(1024, 3)))
	} else if memoryStr[len(memoryStr)-2:] == "mi" {
		val, err := strconv.ParseFloat(memoryStr[:len(memoryStr)-2], 64)
		if err != nil {
			fmt.Println(err.Error())
			return bytes
		}
		bytes = int64(math.Round(val * math.Pow(1024, 2)))
	} else if memoryStr[len(memoryStr)-1:] == "k" {
		val, err := strconv.ParseInt(memoryStr[:len(memoryStr)-1], 10, 64)
		if err != nil {
			fmt.Println(err.Error())
			return bytes
		}
		bytes = val * 1024
	} else if memoryStr[len(memoryStr)-1:] == "b" {
		val, err := strconv.ParseInt(memoryStr[:len(memoryStr)-1], 10, 64)
		if err != nil {
			fmt.Println(err.Error())
			return bytes
		}
		bytes = val
	} else {
		return bytes
	}
	return bytes
}

func ConvertCpuToBytes(cpuStr string) int64 {
	cpuStr = strings.TrimSpace(cpuStr)
	if len(cpuStr) == 0 {
		return 1 * 1e8 //default 0.1 core
	} else {
		cpuLimit, err := strconv.ParseFloat(cpuStr, 64)
		if err != nil {
			fmt.Println(err.Error())
			return -1
		}
		NanoCPU := (int64)(cpuLimit * 1e9)
		return NanoCPU
	}
}
