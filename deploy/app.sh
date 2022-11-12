#!/bin/bash
# - app.sh
#

function start() {
    ./App
}

function stop() {
    pid=`ps -ef | awk '/App/ && !/awk/ && !/app.sh/{print $2}'`
    kill -9 $pid
}

function restart() {
    stop
    sleep 3
    start
}

case $1 in
    start )
    start
        ;;
    stop)
    stop
        ;;
    restart )
    restart
        ;;
    * )

    echo "choose from: [start | stop | restart]"
    ;;
esac
