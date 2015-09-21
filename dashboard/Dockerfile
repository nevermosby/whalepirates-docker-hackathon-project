FROM node:0.12.7
MAINTAINER David Li <wenquan.li@hp.com>

RUN npm install express handlebars jade body-parser request node-etcd

VOLUME ["/usr/src/app"]

EXPOSE 4000

WORKDIR /usr/src/app

CMD ["node", "app.js"]
