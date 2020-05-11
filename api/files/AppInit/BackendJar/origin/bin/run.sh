#!/bin/bash

source /etc/profile
java -jar $JAVA_OPTS $JAR_BOOTER > $APP_LOG 2>&1 &
echo $! > $APP_PID