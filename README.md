# Goal for this project
Hackathon Project for DOCKER GLOBAL HACK DAY #3

# What's this project about
This project aims to create a configuration service based on [etcd][1], which is a distributed consistent key-value store for shared configuration and service discovery, with several features:
- Simple web UI
- Sensitive data encryption
- Etcd data backup with Openstack Swift and restore

# Demo
We will make docker registry as the demo to subscribe etcd configuration, which would greatly help to configure and deploy registry distribution in a complex and varied way. Our solution will initiate or restart registered docker registries when detecting a change from etcd data. In addition, security issues addressing etcd data will be discussed and solved.

# How to run
An etcd watcher will watch for etcd data changes, and update the container that watcher maintains. 
Run a registry
```
docker run -d -p 5000:5000 --restart=always -v /path/to/config.yml:/etc/docker/registry/config.yml --name registry distribution/registry
```
To start a etcd watcher, run
```
docker run -d -v /path/to/config.yml:/config.yml -v /path/to/id_rsa:/id_rsa -v /var/run/docker.sock:/var/run/docker.sock -e ETCD_HOST=<ip>:<port> -e WATCH_ETCD_DIR=hack -e CONFIG=/config.yml -e PRIVATE_KEY=/id_rsa -e SERVE_CONTAINER=<registry container id> --name watcher watcher-image
```
A demo is shown to modify swiftbackend password by setting key=hack/storage/swift/password from web UI we provided. 

Give a test to show new config changes take effect.

# Benefits
You can have configuration data stored in different directory for different releases. For example, etcd directory registry-test/ is for testing, and registry-production-1.0/ is for release 1.0

You can have configuration data stored in different directory for different teams. Each team might use different credentials for swift as storage backends. In that case, we could seperate config data for each team.  

The configuration data never get lost and the history of records will be kept in case you need to restore your services.

Simplify configuration process by UI provided, developer will see clearer and manage eaiser. And automate the deploy to save time.

[1]: https://github.com/coreos/etcd
