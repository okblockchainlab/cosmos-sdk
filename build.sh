#!/usr/bin/env bash

export GO111MODULE=on

export GOPROXY=direct,https://athens.azurefd.net,http://goproxy.io,http://mirrors.aliyun.com/goproxy,https://gocenter.io,


go mod tidy
go mod vendor

#export GO111MODULE=off
#make install