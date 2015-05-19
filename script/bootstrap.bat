@echo off
set GOPATH=%~dp0
echo %GOPATH%

go get -u github.com/golang/protobuf/protoc-gen-go
