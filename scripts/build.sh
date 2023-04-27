#!/bin/bash
export GOPROXY=https://mirrors.aliyun.com/goproxy/
sudo chmod -R 777 .
go build -o kubectl