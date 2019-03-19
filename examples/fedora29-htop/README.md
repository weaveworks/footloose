# Customize the OS image

It is possible to create docker images that specialize a [`footloose` base
image](https://github.com/weaveworks/footloose#choosing-the-os-image-to-run) to
suit your needs.

For instance, if we want the created machines to run `fedora29` with the
`htop` package already pre-installed:

```Dockerfile
FROM quay.io/footloose/fedora29

# Pre-seed the htop package
RUN dnf -y install htop && dnf clean all

```

Build that image:

```console
docker build -t fedora29-htop .
```

Configure `footloose.yaml` to use that image by either editing the file or running:

```console
footloose config create --image fedora29-htop
````

`htop` will be available on the newly created machines!

```console
$ footloose create
$ footloose ssh root@node0
# htop
```
