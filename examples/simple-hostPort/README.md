# Simple port mapping example

First prepare your deploy setup. Notice the last two lines, which means that
node0 port 22 will get mapped to host port 2222, and node1 port 22 will get
mapped to host port 2223:

```console
$ cat footloose.yaml
cluster:
  name: cluster
  privateKey: cluster-key
machines:
- count: 2
  spec:
    image: quay.io/footloose/centos7
    name: node%d
    portMappings:
    - containerPort: 22
      hostPort: 2222
```

Now you can deploy your cluster:

```console
$ footloose create
INFO[0000] Image: quay.io/footloose/centos7 present locally 
INFO[0000] Creating machine: cluster-node0 ...          
INFO[0001] Creating machine: cluster-node1 ...          

```

You now have two container running, listening on SSH port 2222 and 2223 of the host:


```console
$ ssh root@127.0.0.1 -p 2222 -i cluster-key hostname
The authenticity of host '[127.0.0.1]:2222 ([127.0.0.1]:2222)' can't be established.
ECDSA key fingerprint is SHA256:rUXnIB9Nmpy8bzEOcr2MWLVOdkzs9dLSXh7mfP/v7Po.
Are you sure you want to continue connecting (yes/no)? yes
Warning: Permanently added '[127.0.0.1]:2222' (ECDSA) to the list of known hosts.
node0

$ ssh root@127.0.0.1 -p 2223 -i cluster-key hostname
The authenticity of host '[127.0.0.1]:2223 ([127.0.0.1]:2223)' can't be established.
ECDSA key fingerprint is SHA256:0vFd0G655FY1PA/04MZKbT/4dmxP8O+hrzMJs/83uaw.
Are you sure you want to continue connecting (yes/no)? yes
Warning: Permanently added '[127.0.0.1]:2223' (ECDSA) to the list of known hosts.
node1
```

When finished, clean up:

```console
$ footloose delete
INFO[0000] Deleting machine: cluster-node0 ...          
INFO[0000] Deleting machine: cluster-node1 ...      
```

