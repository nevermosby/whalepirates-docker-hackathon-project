FROM ubuntu:14.04
MAINTAINER David Li <wenquan.li@hpe.com>

ENV DEBIAN_FRONTEND noninteractive
ENV ETCD_IP <Your Etcd Server Address>
ENV ETCD_PORT <Your Etcd Server IP>
ENV SWIFT_OBJECT etcd-data-backup.json

RUN echo "deb http://archive.ubuntu.com/ubuntu trusty-backports main restricted universe multiverse" >> /etc/apt/sources.list \
    && apt-get update \
    && apt-get upgrade -yq \
    && apt-get install -yq --no-install-recommends git jq/trusty-backports nodejs npm \
		python-dev build-essential \
		libffi-dev libssl-dev zlib1g-dev libxml2-dev \
		libxslt-dev python-pip \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/* \
    && rm -rf /tmp/* \
    && ln -s /usr/bin/nodejs /usr/bin/node \
    && pip install wrapt \
    && pip install python-keystoneclient==1.6.0 \
    && pip install python-swiftclient==2.4.0

RUN npm install -g git+https://github.com/tombburnell/etcd-dump.git

ADD start.sh /start.sh

#["/start.sh"]
