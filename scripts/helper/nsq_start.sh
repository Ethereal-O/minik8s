nsqlookupd > nsqlookupd.log 2>&1 &
nsqd --lookupd-tcp-address=127.0.0.1:4160 > nsqd.log 2>&1 &
nsqadmin --lookupd-http-address=127.0.0.1:4161 > nsqadmin.log 2>&1 &