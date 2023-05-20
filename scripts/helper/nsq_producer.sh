#!/bin/bash
/bin/nsqlookupd &
sleep 1
/bin/nsqadmin --lookupd-http-address=127.0.0.1:4161 &
sleep 1
/bin/nsqd --lookupd-tcp-address=127.0.0.1:4160