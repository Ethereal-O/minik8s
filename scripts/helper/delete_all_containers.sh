#!/bin/bash
sudo docker kill $(docker ps -q) > /dev/null 2>&1
sudo docker rm $(docker ps -aq) > /dev/null 2>&1