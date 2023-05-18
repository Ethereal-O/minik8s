#!/bin/bash
export GOPROXY=https://goproxy.io,direct
export GO111MODULE=on
sudo chmod -R 777 .
go build -o kubectl