# Running manually with Go

## Requirements
* Go >= v1.6
* libpcap >= 1.6.2

## Run
* execute ./run-sh

```sh
$ ./run-sh
```

## Example
[This Dockerfile](https://github.com/siaflab/haskap-jam-pack/blob/master/server/log-interceptor/Dockerfile) describes how to run haskap-jam-interceptor on linux_amd64 platform.

```sh
$ docker build .
Sending build context to Docker daemon 7.217 MB
Step 1 : FROM golang:1.6
 ---> bf9e64d14fb2
Step 2 : RUN apt-get update 	&& apt-get -y install libpcap-dev 	&& rm -rf /var/cache/apt/archives/*.deb
 ---> Running in 1f1941abe63c
Get:1 http://security.debian.org jessie/updates InRelease [63.1 kB]
Get:2 http://security.debian.org jessie/updates/main amd64 Packages [269 kB]
Ign http://httpredir.debian.org jessie InRelease
Get:3 http://httpredir.debian.org jessie-updates InRelease [136 kB]
Get:4 http://httpredir.debian.org jessie Release.gpg [2373 B]
Get:5 http://httpredir.debian.org jessie Release [148 kB]
Get:6 http://httpredir.debian.org jessie-updates/main amd64 Packages [3619 B]
Get:7 http://httpredir.debian.org jessie/main amd64 Packages [9034 kB]
Fetched 9656 kB in 26s (362 kB/s)
Reading package lists...
Reading package lists...
Building dependency tree...
Reading state information...
The following extra packages will be installed:
  libpcap0.8 libpcap0.8-dev
The following NEW packages will be installed:
  libpcap-dev libpcap0.8 libpcap0.8-dev
0 upgraded, 3 newly installed, 0 to remove and 1 not upgraded.
Need to get 385 kB of archives.
After this operation, 985 kB of additional disk space will be used.
Get:1 http://httpredir.debian.org/debian/ jessie/main libpcap0.8 amd64 1.6.2-2 [133 kB]
Get:2 http://httpredir.debian.org/debian/ jessie/main libpcap0.8-dev amd64 1.6.2-2 [229 kB]
Get:3 http://httpredir.debian.org/debian/ jessie/main libpcap-dev all 1.6.2-2 [23.6 kB]
debconf: delaying package configuration, since apt-utils is not installed
Fetched 385 kB in 2s (142 kB/s)
Selecting previously unselected package libpcap0.8:amd64.
(Reading database ... 14702 files and directories currently installed.)
Preparing to unpack .../libpcap0.8_1.6.2-2_amd64.deb ...
Unpacking libpcap0.8:amd64 (1.6.2-2) ...
Selecting previously unselected package libpcap0.8-dev.
Preparing to unpack .../libpcap0.8-dev_1.6.2-2_amd64.deb ...
Unpacking libpcap0.8-dev (1.6.2-2) ...
Selecting previously unselected package libpcap-dev.
Preparing to unpack .../libpcap-dev_1.6.2-2_all.deb ...
Unpacking libpcap-dev (1.6.2-2) ...
Setting up libpcap0.8:amd64 (1.6.2-2) ...
Setting up libpcap0.8-dev (1.6.2-2) ...
Setting up libpcap-dev (1.6.2-2) ...
Processing triggers for libc-bin (2.19-18+deb8u3) ...
 ---> 7ced4eba948e
Removing intermediate container 1f1941abe63c
Step 3 : ADD . /home
 ---> 93b741c1d120
Removing intermediate container 0831bd3bd863
Step 4 : WORKDIR /home
 ---> Running in 152063a4e313
 ---> 83259da6ea45
Removing intermediate container 152063a4e313
Step 5 : RUN ip a
 ---> Running in 4a5de850ba71
1: lo: <LOOPBACK,UP,LOWER_UP> mtu 65536 qdisc noqueue state UNKNOWN group default
    link/loopback 00:00:00:00:00:00 brd 00:00:00:00:00:00
    inet 127.0.0.1/8 scope host lo
       valid_lft forever preferred_lft forever
    inet6 ::1/128 scope host
       valid_lft forever preferred_lft forever
47: eth0@if48: <BROADCAST,MULTICAST,UP,LOWER_UP> mtu 1500 qdisc noqueue state UP group default
    link/ether 02:42:ac:11:00:02 brd ff:ff:ff:ff:ff:ff
    inet 172.17.0.2/16 scope global eth0
       valid_lft forever preferred_lft forever
    inet6 fe80::42:acff:fe11:2/64 scope link tentative
       valid_lft forever preferred_lft forever
 ---> 70170729915d
Removing intermediate container 4a5de850ba71
Step 6 : RUN sed -i -e "s/\"deviceName\": \"lo0\"/\"deviceName\": \"lo\"/g" haskap-jam-interceptor-config.json
 ---> Running in f2ac7f438a66
 ---> ff3776849e30
Removing intermediate container f2ac7f438a66
Step 7 : RUN go get github.com/google/gopacket; go run haskap-jam-interceptor.go
 ---> Running in 728130f7c292
config.DeviceName: lo
config.ReceivePort: 4558
config.SendToAddress: 127.0.0.1
config.SendToPort: 3333
#####
2016-02-29 05:04:59.552581323 +0000 UTC
haskap-jam-interceptor started successfully.
version: , build: , date:
capturing UDP port 4558 packets.
and will send to 127.0.0.1:3333
#####
```
