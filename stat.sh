#!/bin/bash
pid=$(ps -ef | awk '$8=="mongod" && $9=="--shardsvr" {print $2}')
curTime=$(date +%Y%m%d%H%M%S)
if [ "$1" == "stop" ]; then
    killall -9 ifstat
    killall -9 psrecord
    killall -9 iostat
    killall -9 mongotop
elif [ "$1" == "start" ]; then
    while [[ $# -gt 0 ]]; do
        key="$2"
        case "$key" in
            -d|--delay)
            shift
            delay=$2
            ;;
            -h|--mongo-host)
            shift
            host=$2
            ;;
            -p|--mongo-port)
            shift
            port=$2
            ;;
            *)
            if [ "$key" != "" ]; then
                echo "Unknown option '$key'"
            fi
            ;;
        esac
        shift
done

    if [ "$delay" == "" ]; then
        delay=1
    fi
    if [ "$port" == "" ]; then
	port=27018
    fi
    echo "Starting Stats Collection at" $curTime "with delay" $delay
    if [ "$pid" == "" ]; then
        echo "Could not find 'mongod --shardsvr' running"
    else
        psrecord $pid --interval $delay --log $curTime-$pid".psrecord" &
    fi
    ifstat $delay > $curTime".ifstat" &
    iostat $delay > $curTime".iostat" &
    if [ "$host" != "" ]; then
    	mongotop --host $host --port $port $delay &
        mongostat --host $host --port $port $delay > $curTime".mongostat" &
    fi
else
    echo "Usage: ./stat.sh start|stop [-d|--delay <delay>] [-h|--mongo-host <XXX.XXX.XXX.XXX>] [-p|--mongo-port <port>"
fi
