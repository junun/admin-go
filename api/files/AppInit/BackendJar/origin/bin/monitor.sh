#!/bin/bash
#
#

SHELL_DIR=$(cd $(dirname $0);pwd)

LOCK_FILE=/dev/shm/`echo ${SHELL_DIR}|sed 's!/!.!g;s!.!!'`.monitor.lock

{
    flock -n 100 || { exit 2; }

    cd ${SHELL_DIR}
    function monitor() {
        while true;do
            ./control monitor
            sleep 3
        done
    }

    monitor >> ../logs/monitor.log 2>&1 &

} 100<>${LOCK_FILE}