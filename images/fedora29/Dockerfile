FROM fedora:29

ENV container docker

RUN dnf -y install sudo openssh-server procps-ng hostname net-tools iproute iputils wget && dnf clean all

EXPOSE 22

# https://www.freedesktop.org/wiki/Software/systemd/ContainerInterface/
STOPSIGNAL SIGRTMIN+3

CMD ["/bin/bash"]
