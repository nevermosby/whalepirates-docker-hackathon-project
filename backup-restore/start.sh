#!/bin/sh

MODE=
KEY=

if [ "$1" == "backup" ]; then
    MODE='dump'
elif [ "$1" == "restore" ]; then
    MODE='restore'
else
    >&2 echo "You must provide either backup or restore to run this container"
    exit 64
fi

if [ -z "$2" ]; then
    KEY=/
else
    KEY=$2
fi

if [ "$MODE" == "restore" ]; then
    if [ -f "/tmp/dump.json" ]; then
        echo "dump.json already provided, skipping download"
    else
        swift --insecure download $SWIFT_CONTAINER $SWIFT_OBJECT

        if [ $? != 0 ]; then
            >&2 echo "There was a problem fetching the backup from Openstack Swift"
            exit $?
        fi
    fi

    jq '[.[] | select(.key | startswith("'"$KEY"'"))]' /tmp/dump.json > /tmp/tmp.json
    mv /tmp/tmp.json /tmp/dump.json
fi

etcd-dump -h $ETCD_IP -p $ETCD_PORT -f /tmp/dump.json $MODE

if [ $? != 0 ]; then
    exit $?
fi

if [ "$MODE" == "dump" ]; then
    swift --insecure upload $SWIFT_CONTAINER $SWIFT_OBJECT

    if [ $? != 0 ]; then
        exit $?
    fi
fi

exit 0

