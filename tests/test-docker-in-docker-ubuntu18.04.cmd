footloose create --config %t.yaml
footloose --config %t.yaml ssh root@node0 -- apt-get update && apt-get install -y docker.io
footloose --config %t.yaml ssh root@node0 systemctl start docker
footloose --config %t.yaml ssh root@node0 docker pull busybox
%out footloose --config %t.yaml ssh root@node0 docker run busybox echo success
footloose delete --config %t.yaml
