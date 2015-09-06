#!/bin/bash

set -e
set -x

sudo apt-get update

curl -sSL https://get.docker.com/ | sh

sudo usermod -aG docker $USER
newgrp docker

docker pull ubuntu:trusty

cd $HOME
mkdir docker
cd docker

cat <<EOF > Dockerfile
FROM ubuntu:trusty

RUN apt-get update && apt-get install -y \
  openjdk-6-jre \
  openjdk-7-jre \
  ruby2.0 \
  python \
  python3

RUN sudo useradd -U contestant

USER contestant
EOF

docker build -t combo .
