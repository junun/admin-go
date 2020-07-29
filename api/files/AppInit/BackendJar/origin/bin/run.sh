#!/bin/bash

source /etc/profile

WORKSPACE=$(cd $(dirname $0)/; pwd)
cd $WORKSPACE
source var
app=${AppName}
JAVA_OPTS=${JavaOpts}
JAR_BOOTER="${app}.jar"
APP_BASE=/data/webapps/$app
APP_PID=${APP_BASE}/${app}.pid
APP_LOG=$APP_BASE/logs/${app}.log

cd $(dirname $(find -L $APP_BASE -name $JAR_BOOTER | head -1))

java -jar ${JAVA_OPTS} ${JAR_BOOTER} > ${APP_LOG} 2>&1 &
echo $! > ${APP_PID}