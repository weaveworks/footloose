FROM quay.io/footloose/ubuntu18.04

RUN apt-get update && apt-get install -y apache2
COPY index.html /var/www/html

RUN systemctl enable apache2.service

EXPOSE 80
