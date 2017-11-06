#!/bin/bash

#description: 这是一个系统服务脚本
#chkconfig: 2345 20 81 

EXEC_PATH=/home/macd/blogci/bin          #需要注册的可执行文件所在文件夹的路径
EXEC=blogci                  #需要注册的可执行文件的文件名
DAEMON="$EXEC_PATH/$EXEC"     #需要注册的可执行文件的完整文件名
PID_FILE=/var/run/blogci.pid #指定运行时的进程号

. /etc/rc.d/init.d/functions   #调用通用方法

if [ ! -x $DAEMON ]; then
    echo "ERROR: $DAEMON not found"
    exit 1
fi

stop()
{
    echo "Stoping $EXEC ..."
    ps aux | grep "$DAEMON" | kill -9 `awk '{print $2}'` >/dev/null 2>&1
    rm -f $PID_FILE
    usleep 100
    echo "Shutting down $EXEC: [  OK  ]"
}

start()
{
    echo "Starting $EXEC ..."
    $DAEMON > /dev/null &
    pidof $EXEC > $PID_FILE
    usleep 100
    echo "Starting $EXEC: [  OK  ]"
}

restart()
{
    stop
    start
}

case "$1" in
    start)
        start
        ;;
    stop)
        stop
        ;;
    restart)
        restart
        ;;
    status)
        status -p $PID_FILE $DAEMON
        ;;
    *)
        echo "Usage: service $EXEC {start|stop|restart|status}"
        exit 1
esac

exit $?
