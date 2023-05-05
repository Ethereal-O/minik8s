package fileServer

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"minik8s/pkg/client"
	"minik8s/pkg/object"
	"minik8s/pkg/util/config"
	"os"
	"path/filepath"
)

func UploadFile(dirpath string, key string) {
	var files []string
	_ = filepath.Walk(dirpath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			filename := info.Name()
			data, err := ioutil.ReadFile(path)
			if err != nil {
				return err
			}
			files = append(files, filename)
			files = append(files, string(data))
		}
		return nil
	})
	filesJson, err := json.Marshal(files)
	if err != nil {
		fmt.Println(err.Error())
	}
	gpufile := object.GpuFile{Dirname: key, Data: string(filesJson)}
	gpufileJson, err := json.Marshal(gpufile)
	if err != nil {
		fmt.Println(err.Error())
	}
	client.Put_object(key, string(gpufileJson), config.GPUFILE_TYPE)
}

func DownloadFile(value string, dirpath string) {
	//Step1: Create the dir in the shareDir
	err := os.MkdirAll(dirpath, 0777)
	if err != nil {
		return
	}

	var files []string
	err = json.Unmarshal([]byte(value), &files)
	if err != nil {
		return
	}
	var flag = 0
	var filename string
	for _, item := range files {
		if flag%2 == 0 {
			//Step2: Create the file in the dir
			filename = item
			fmt.Println(filename)
			os.Create(dirpath + "/" + filename)
		} else {
			//Step3: Write the data to the file
			file, err := os.OpenFile(dirpath+"/"+filename, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0777)
			if err != nil {
				fmt.Println(err)
				return
			}
			size, err := file.Write([]byte(item))
			fmt.Printf("The %s size is %d bytes\n", filename, size)
		}
		flag++
	}
}
