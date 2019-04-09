FROM opensuse/tumbleweed:latest

RUN zypper -n install systemd; zypper clean ; \
(cd /usr/lib/systemd/system/sysinit.target.wants/; for i in *; do [ $i == systemd-tmpfiles-setup.service ] || rm -f $i; done); \
rm -f /usr/lib/systemd/system/multi-user.target.wants/*;\
rm -f /etc/systemd/system/*.wants/*;\
rm -f /usr/lib/systemd/system/local-fs.target.wants/*; \
rm -f /usr/lib/systemd/system/sockets.target.wants/*udev*; \
rm -f /usr/lib/systemd/system/sockets.target.wants/*initctl*; \
rm -f /usr/lib/systemd/system/basic.target.wants/*;\
rm -f /usr/lib/systemd/system/anaconda.target.wants/*;

RUN zypper --non-interactive install \
  openssh \
  sudo hostname iproute2 net-tools iputils wget && \
  zypper clean

RUN ln -s /bin/systemd /sbin/init
EXPOSE 22

# https://www.freedesktop.org/wiki/Software/systemd/ContainerInterface/
STOPSIGNAL SIGRTMIN+3

CMD ["bash"]