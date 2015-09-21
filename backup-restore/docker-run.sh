#!/bin/sh

docker run \
	-e OS_AUTH_URL= \
	-e OS_USERNAME= \
	-e OS_PASSWORD="" \
	-e OS_TENANT_ID= \
	-e OS_AUTH_VERSION=2 \
	-e ETCD_IP= \
	-e ETCD_PORT=4001 \
	-e SWIFT_CONTAINER= \
	david/etcd-backup-restore:0.1 bash /start.sh backup
