#!/bin/bash
sh ./scripts/clean.sh

killall nsqlookupd > /dev/null 2>&1
killall nsqd > /dev/null 2>&1
killall nsqadmin > /dev/null 2>&1
echo "[Stop] NSQ Producer/Consumer stopped!" 1>&2
weave reset > /dev/null 2>&1
echo "[Stop] Weave subnet stopped!" 1>&2