#!/bin/bash

GOPATH=`pwd`
#if [ ! -f "$GOPATH/bin/protoc-gen-go" ]; then  
#    GOPATH=$GOPATH go install code.google.com/p/goprotobuf/protoc-gen-go
#fi  

protoc --go_out=src protocol/*.proto
protoc --python_out=emulator protocol/*.proto
