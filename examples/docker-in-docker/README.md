# Running `dockerd` in Container Machines

To run `dockerd` inside a docker container, two things are needed:

- Run the container as privileged (we could probably do better! expose
capabilities instead).
- Mount `/var/lib/docker` as volume, here an anonymous volume. This is
because of [limitations][dind] of what you can do with the overlay system
docker is setup to use.

```yaml
cluster:
  name: cluster
  privateKey: cluster-key
machines:
- count: 1
  spec:
    image: quay.io/footloose/centos7
    name: node%d
    portMappings:
    - containerPort: 22
    privileged: true
    volumes:
    - type: volume
      destination: /var/lib/docker
```

You can then install and run docker on the machine:

```console
$ footloose create
$ footloose ssh root@node0
# yum install -y docker iptables
[...]
# systemctl start docker
# docker run busybox echo 'Hello, World!'
Hello, World!
```

[dind]: https://jpetazzo.github.io/2015/09/03/do-not-use-docker-in-docker-for-ci/
