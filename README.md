[![Build Status](https://travis-ci.com/dlespiau/footloose.svg?branch=master)](https://travis-ci.com/dlespiau/footloose)

# footloose

`footloose` is a developer tool creating containers that look like virtual
machines. Those containers run `systemd` as PID 1 and a ssh daemon that can
be used to login into the container. Such "machines" behave very much like a
VM, it's even possible to run `dockerd` in them :)

`footloose` in action:

[![asciicast](https://asciinema.org/a/226185.svg)](https://asciinema.org/a/226185)

`footloose` can be used for a variety of tasks, wherever you'd like virtual
machines but don't want to wait for them to boot or need many of them. An
easy way to think about it is: [Vagrant](https://www.vagrantup.com/), but
with containers.

## Install

`footloose` hasn't reached a first version yet but you can install it from sources:

```
go get github.com/dlespiau/footloose
```

## Usage

`footloose` reads a description of the *Cluster* of *Machines* to create from a
file, by default named `Footloose`. The `config` command helps with creating the
initial config file:

```console
# Create a Footloose config file. Instruct we want to create 3 machines instead
# of the default, 1.
footloose config --replicas 3
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
$ footloose ssh node1
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
