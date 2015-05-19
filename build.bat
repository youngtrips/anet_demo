@echo off
set GOPATH=%~dp0
echo %GOPATH%
go install ./src/cmd/scorpio

