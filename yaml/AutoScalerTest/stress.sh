#!/bin/bash

echo $1

while true
do
  curl -sS http://$1:80 > /dev/null
#  curl -sS http://10.10.1.2:80 > /dev/null
#  curl -sS http://10.10.1.3:80 > /dev/null
done