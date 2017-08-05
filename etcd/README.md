Overview
=====

Build
=====

The project could be build into a Dockerized server with 'make docker-build'. An local image named 'etcd_uid' is available after that.
    
    $make
    $docker build -t etcd_uid .
Run
=====

Start etcd
-----
   docker run -d --name etcd -p 4001:4001 -p 7001:7001 appcelerator/etcd

Run etcd_uid
=====
   $docker run --name container_name -p 9090:9090 etcd_uid
Send curl command
=====
    $curl 192.168.1.175:8080/user/etcduid