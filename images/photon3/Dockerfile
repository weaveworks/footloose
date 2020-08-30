FROM photon:3.0

ENV container docker

RUN yum -y install \
    systemd dbus sudo procps-ng net-tools iproute iputils wget vim openssh-server && \
    yum clean all

RUN echo 'root:*:17995::::::' > /etc/shadow

EXPOSE 22

# https://www.freedesktop.org/wiki/Software/systemd/ContainerInterface/
STOPSIGNAL SIGRTMIN+3

CMD ["/bin/bash"]
