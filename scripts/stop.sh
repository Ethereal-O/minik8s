sh ./scripts/clean.sh

killall nsqlookupd > /dev/null 2>&1
killall nsqd > /dev/null 2>&1
killall nsqadmin > /dev/null 2>&1
weave stop > /dev/null 2>&1