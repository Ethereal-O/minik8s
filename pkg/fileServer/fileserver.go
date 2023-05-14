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
	fileChan, stopFunc := messging.Watch("/"+config.TRANSFILE_TYPE, true)
	dealFile(fileChan)
	fmt.Println("File server start")

	// Wait until Ctrl-C
	<-FileServerToExit
	stopFunc()
	FileServerExited <- true
}

func dealFile(fileChan chan string) {
	for {
		select {
		case mes := <-fileChan:
			if mes == "hello" {
				continue
			}
			var file object.TransFile
			json.Unmarshal([]byte(mes), &file)
			filename := file.Dirname
			filedata := file.Data
			fmt.Println("[downloader receive]" + filename)
			if file.Tp == "GpuFile" {
				DownloadFile(filedata, config.GPU_NODE_DIR_PATH+"/"+filename)
			}
			if file.Tp == "FuncFile" {
				DownloadFile(filedata, config.FUNC_NODE_DIR_PATH+"/"+filename)
			}
		}
	}
}
