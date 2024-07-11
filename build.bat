@echo off

set GOOS=windows
go build -ldflags="-s -w" -trimpath -o server_windows.exe

set GOOS=linux
go build -ldflags="-s -w" -trimpath -o server