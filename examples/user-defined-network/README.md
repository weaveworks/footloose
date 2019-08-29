# Using user-defined network example

Using a user-defined network enables DNS name resolution of the container names, so you can talk
to each container of the cluster just using the hostname.

First prepare your deploy setup. Notice the line 'network' which specifies which user-defined network the containers should be attached to.

```console
$ cat footloose.yaml
cluster:
  name: cluster
  privateKey: cluster-key
machines:
- count: 3
  spec:
    image: quay.io/footloose/centos7:0.6.0
    name: node%d
    network: footloose-cluster
    portMappings:
    - containerPort: 22
```

The user-defined network has to be created manually before deploying your cluster:

```console
$ docker network create footloose-cluster
c558b7218393a2e4c89b19f7904d244192664997f46eb6edfc3217e187472afc
```

Now you can deploy your cluster:

```console
$ footloose create
INFO[0000] Image: quay.io/footloose/centos7 present locally
INFO[0000] Creating machine: cluster-node0 ...
INFO[0001] Creating machine: cluster-node1 ...
INFO[0002] Creating machine: cluster-node2 ...

```

You now have three containers running, which can talk to each other using their hostnames:

```console
$ footloose ssh root@node0
[root@node0 ~]# ping -c 4 node1
PING node1 (172.25.0.3) 56(84) bytes of data.
64 bytes from cluster-node1.footloose-cluster (172.25.0.3): icmp_seq=1 ttl=64 time=0.240 ms
64 bytes from cluster-node1.footloose-cluster (172.25.0.3): icmp_seq=2 ttl=64 time=0.289 ms
64 bytes from cluster-node1.footloose-cluster (172.25.0.3): icmp_seq=3 ttl=64 time=0.193 ms
64 bytes from cluster-node1.footloose-cluster (172.25.0.3): icmp_seq=4 ttl=64 time=0.205 ms

--- node1 ping statistics ---
4 packets transmitted, 4 received, 0% packet loss, time 3044ms
rtt min/avg/max/mdev = 0.193/0.231/0.289/0.041 ms
[root@node0 ~]# ping -c 4 node2
PING node2 (172.25.0.4) 56(84) bytes of data.
64 bytes from cluster-node2.footloose-cluster (172.25.0.4): icmp_seq=1 ttl=64 time=0.109 ms
64 bytes from cluster-node2.footloose-cluster (172.25.0.4): icmp_seq=2 ttl=64 time=0.184 ms
64 bytes from cluster-node2.footloose-cluster (172.25.0.4): icmp_seq=3 ttl=64 time=0.143 ms

--- node2 ping statistics ---
3 packets transmitted, 3 received, 0% packet loss, time 2059ms
rtt min/avg/max/mdev = 0.109/0.145/0.184/0.032 ms

```

When finished, clean up:

```console
$ footloose delete
INFO[0000] Deleting machine: cluster-node0 ...
INFO[0000] Deleting machine: cluster-node1 ...
INFO[0001] Deleting machine: cluster-node2 ...
```

