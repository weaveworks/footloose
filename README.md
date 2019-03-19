[![Build Status](https://travis-ci.org/weaveworks/footloose.svg?branch=master)](https://travis-ci.org/weaveworks/footloose)
[![Go Report Card](https://goreportcard.com/badge/github.com/weaveworks/footloose)](https://goreportcard.com/report/github.com/weaveworks/footloose)
[![GoDoc](https://godoc.org/github.com/weaveworks/footloose?status.svg)](https://godoc.org/github.com/weaveworks/footloose)

# footloose

`footloose` is a developer tool creating containers that look like virtual
machines. Those containers run `systemd` as PID 1 and a ssh daemon that can
be used to login into the container. Such "machines" behave very much like a
VM, it's even possible to run [`dockerd` in them][readme-did] :)

`footloose` can be used for a variety of tasks, wherever you'd like virtual
machines but don't want to wait for them to boot or need many of them. An
easy way to think about it is: [Vagrant](https://www.vagrantup.com/), but
with containers.

`footloose` in action:

[![asciicast](https://asciinema.org/a/226185.svg)](https://asciinema.org/a/226185)

[readme-did]: https://github.com/weaveworks/footloose#running-dockerd-in-container-machines

## Install

`footloose` binaries can be downloaded from the [release page][gh-release]:

```console
# For Linux
curl -Lo footloose https://github.com/weaveworks/footloose/releases/download/0.2.0/footloose-0.2.0-linux-x86_64
chmod +x footloose
sudo mv footloose /usr/local/bin/

# For macOS
curl -Lo footloose https://github.com/weaveworks/footloose/releases/download/0.2.0/footloose-0.2.0-darwin-x86_64
chmod +x footloose
sudo mv footloose /usr/local/bin/
```

Alternatively, build and install `footloose` from sources. It requires having
`go >= 1.11` installed:

```console
GO111MODULE=on go get github.com/weaveworks/footloose
```

[gh-release]: https://github.com/weaveworks/footloose/releases

## Usage

`footloose` reads a description of the *Cluster* of *Machines* to create from a
file, by default named `Footloose`. The `config` command helps with creating the
initial config file:

```console
# Create a Footloose config file. Instruct we want to create 3 machines instead
# of the default, 1.
footloose config create --replicas 3
```

Start the cluster of 3 machines:

```console
$ footloose create
INFO[0000] Pulling image: quay.io/footloose/centos7 ...
INFO[0007] Creating machine: cluster-node0 ...
INFO[0008] Creating machine: cluster-node1 ...
INFO[0008] Creating machine: cluster-node2 ...
```

> It only takes a second to create those machines. The first time `create`
runs, it will pull the docker image used by the `footloose` containers so it
will take a tiny bit longer.

SSH into a machine with:

```console
$ footloose ssh root@node1
[root@1665288855f6 ~]# ps fx
  PID TTY      STAT   TIME COMMAND
    1 ?        Ss     0:00 /sbin/init
   23 ?        Ss     0:00 /usr/lib/systemd/systemd-journald
   58 ?        Ss     0:00 /usr/sbin/sshd -D
   59 ?        Ss     0:00  \_ sshd: root@pts/1
   63 pts/1    Ss     0:00      \_ -bash
   82 pts/1    R+     0:00          \_ ps fx
   62 ?        Ss     0:00 /usr/lib/systemd/systemd-logind
```

## Choosing the OS image to run

`footloose` will default to running a centos 7 container image. The `--image`
argument of `config create` can be used to configure the OS image. Valid OS
images are:

- `quay.io/footloose/centos7`
- `quay.io/footloose/fedora29`
- `quay.io/footloose/ubuntu16.04`
- `quay.io/footloose/ubuntu18.04`
- `quay.io/footloose/amazonlinux2`
- `quay.io/footloose/debian10`

For example:

```console
footloose config create --replicas 3 --image quay.io/footloose/fedora29
```

## `footloose.yaml`

`footloose config create` creates a `footloose.yaml` configuration file that is then
used by subsequent commands such as `create`, `delete` or `ssh`. If desired,
the configuration file can be named differently and supplied with the
`-c, --config` option.

```console
$ footloose config create --replicas 3
$ cat footloose.yaml
cluster:
  name: cluster
  privateKey: cluster-key
machines:
- count: 3
  spec:
    image: quay.io/footloose/centos7
    name: node%d
    portMappings:
    - containerPort: 22
```

This configuration can naturally be edited by hand. The full list of
available parameters are in [the reference documentation][pkg-config].

[pkg-config]: https://godoc.org/github.com/weaveworks/footloose/pkg/config

## Examples

Interesting things can be done with `footloose`!

* [Customize the OS image](./examples/fedora29-htop/README.md)
* [Ansible example](./examples/ansible/README.md)
* [Run Apache in a footloose machine](./examples/apache/README.md)

## Running `dockerd` in Container Machines

To run `dockerd` inside a docker container, two things are needed:

- Run the container as privileged (we could probably do better! expose
capabilities instead!).
- Mount `/var/lib/docker` as volume, here an anonymous volume. This is
because of a [limitations][dind] of what you can do with the overlay systems
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

## Under the hood

Under the hood, *Container Machines* are just containers. They can be
inspected with `docker`:

```console
$ docker ps
CONTAINER ID    IMAGE                        COMMAND         NAMES
04c27967f76e    quay.io/footloose/centos7    "/sbin/init"    cluster-node2
1665288855f6    quay.io/footloose/centos7    "/sbin/init"    cluster-node1
5134f80b733e    quay.io/footloose/centos7    "/sbin/init"    cluster-node0
```

The container names are derived from `cluster.name` and
`cluster.machines[].name`.

They run `systemd` as PID 1, it's even possible to inspect the boot messages:

```console
$ docker logs cluster-node1
systemd 219 running in system mode.
Detected virtualization docker.
Detected architecture x86-64.

Welcome to CentOS Linux 7 (Core)!

Set hostname to <1665288855f6>.
Initializing machine ID from random generator.
Failed to install release agent, ignoring: File exists
[  OK  ] Created slice Root Slice.
[  OK  ] Created slice System Slice.
[  OK  ] Reached target Slices.
[  OK  ] Listening on Journal Socket.
[  OK  ] Reached target Local File Systems.
         Starting Create Volatile Files and Directories...
[  OK  ] Listening on Delayed Shutdown Socket.
[  OK  ] Reached target Swap.
[  OK  ] Reached target Paths.
         Starting Journal Service...
[  OK  ] Started Create Volatile Files and Directories.
[  OK  ] Started Journal Service.
[  OK  ] Reached target System Initialization.
[  OK  ] Started Daily Cleanup of Temporary Directories.
[  OK  ] Reached target Timers.
[  OK  ] Listening on D-Bus System Message Bus Socket.
[  OK  ] Reached target Sockets.
[  OK  ] Reached target Basic System.
         Starting OpenSSH Server Key Generation...
         Starting Cleanup of Temporary Directories...
[  OK  ] Started Cleanup of Temporary Directories.
[  OK  ] Started OpenSSH Server Key Generation.
         Starting OpenSSH server daemon...
[  OK  ] Started OpenSSH server daemon.
[  OK  ] Reached target Multi-User System.
```

## FAQ

### Is `footloose` just like LXD?
In principle yes, but it will also work with Docker container images and
on MacOS as well.

## Help

We are a very friendly community and love questions, help and feedback.

If you have any questions, feedback, or problems with `footloose`:

- Check out the [examples](examples).
- Join the discussion
  - Invite yourself to the <a href="https://slack.weave.works/" target="_blank">Weave community</a> Slack.
  - Ask a question on the [#general](https://weave-community.slack.com/messages/general/) Slack channel.
  - Join the [Weave User Group](https://www.meetup.com/pro/Weave/) and get invited to online talks, hands-on training and meetups in your area.
- [File an issue](https://github.com/weaveworks/footloose/issues/new).

Your feedback is always welcome!
