#!/bin/bash

export GOPATH=`pwd`
echo "GOPATH="$GOPATH

go get github.com/go-sql-driver/mysql
go get github.com/kesselborn/go-getopt
#go get gopkg.in/mgo.v2
go get -u github.com/golang/protobuf/{proto,protoc-gen-go}
