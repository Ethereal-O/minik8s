#!/bin/bash
master_ip=$(cat /master_ip.txt)
/bin/nsqd --lookupd-tcp-address="$master_ip:4160"