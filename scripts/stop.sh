#!/bin/bash
sh ./scripts/clean.sh

sudo killall nsqlookupd > /dev/null 2>&1
sudo killall nsqd > /dev/null 2>&1
sudo killall nsqadmin > /dev/null 2>&1
echo "[Stop] NSQ Producer/Consumer stopped!" 1>&2
