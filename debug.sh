#!/bin/sh

DOCKER=`which docker`

make clean
make
cp build/pgskail.linux docker/pgskail
$DOCKER build -t pires/pgskail docker
$DOCKER run --rm -e ETCD=192.168.59.104 pires/pgskail