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
