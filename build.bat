@echo off

set GOOS=windows
go build -ldflags="-s -w" -trimpath -o server_windows.exe

go build -ldflags="-s -w" -trimpath -o client_windows.exe

set GOOS=linux
go build -ldflags="-s -w" -trimpath -o server

go build -ldflags="-s -w" -trimpath -o client
