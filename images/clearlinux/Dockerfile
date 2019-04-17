FROM clearlinux:latest

ENV container docker

RUN swupd bundle-add openssh-server vim network-basic sudo
RUN echo 'root:*:17995::::::' > /etc/shadow

EXPOSE 22

STOPSIGNAL SIGRTMIN+3

CMD ["/bin/bash"]
