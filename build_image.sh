#!/bin/sh


./go_build_linux_amd64.sh ./build/function-manager

docker build -t 10.10.15.51/fc/function-manager:beta .