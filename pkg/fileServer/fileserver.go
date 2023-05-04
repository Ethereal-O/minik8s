package fileServer

import (
	"encoding/json"
	"fmt"
	"minik8s/pkg/messging"
	"minik8s/pkg/object"
	"minik8s/pkg/util/config"
)

var FileServerExited = make(chan bool)
var FileServerToExit = make(chan bool)

func Start_Fileserver() {
	gpufileChan, stopFunc := messging.Watch("/"+config.GPUFILE_TYPE, true)
	dealFile(gpufileChan)
	fmt.Println("GpuJob Controller start")

	// Wait until Ctrl-C
	<-FileServerToExit
	stopFunc()
	FileServerExited <- true
}

func dealFile(gpufileChan chan string) {
	for {
		select {
		case mes := <-gpufileChan:
			if mes == "hello" {
				continue
			}
			var gpufile object.GpuFile
			json.Unmarshal([]byte(mes), &gpufile)
			jobname := gpufile.Dirname
			jobdata := gpufile.Data
			DownloadFile(jobdata, config.NODE_DIR_PATH+"/"+jobname)
		}
	}
}
