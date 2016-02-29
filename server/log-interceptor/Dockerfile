# an example Dockerfile describes how to run haskap-jam-interceptor.go on linux_amd64

# base image: Go 1.6 on Debian Jessie
FROM golang:1.6

# install libpcap-dev
RUN apt-get update \
	&& apt-get -y install libpcap-dev \
	&& rm -rf /var/cache/apt/archives/*.deb

ADD . /home
WORKDIR /home

RUN ip a
RUN sed -i -e "s/\"deviceName\": \"lo0\"/\"deviceName\": \"lo\"/g" haskap-jam-interceptor-config.json

# already super user and `sudo` not required.
# RUN ./run.sh
RUN go get github.com/google/gopacket; go run haskap-jam-interceptor.go
