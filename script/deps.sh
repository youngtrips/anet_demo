#!/bin/bash

GOPATH=`pwd`
GOPATH=$GOPATH go get -u code.google.com/p/goprotobuf/{proto,protoc-gen-go}
#GOPATH=$GOPATH go get -u code.google.com/p/goprotobuf/...
#go get cjones.org/hg/go-xmpp2.hg/xmpp
