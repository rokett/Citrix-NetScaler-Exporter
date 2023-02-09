@echo off
SETLOCAL

set APP=Citrix-NetScaler-Exporter
set VERSION=4.6.0
set BINARY-WINDOWS-X64=%APP%_%VERSION%_Windows_amd64.exe
set BINARY-LINUX=%APP%_%VERSION%_Linux_amd64

REM Set build number from git commit hash
for /f %%i in ('git rev-parse HEAD') do set BUILD=%%i

set LDFLAGS=-ldflags "-X main.version=%VERSION% -X main.build=%BUILD% -s -w -extldflags '-static'"

goto build

:build
    echo "=== Building docker image ==="
    rem docker build -t rokett/citrix-netscaler-exporter:latest -t rokett/citrix-netscaler-exporter:v%VERSION% .
    rem docker push rokett/citrix-netscaler-exporter:v%VERSION%
    rem docker push rokett/citrix-netscaler-exporter:latest

    echo "=== Building Windows x64 ==="
    set GOOS=windows
    set GOARCH=amd64

    go build -mod=vendor %LDFLAGS% -o %BINARY-WINDOWS-X64%

    echo "=== Building Linux x64 ==="
    set GOOS=linux
    set GOARCH=amd64

    go build -mod=vendor %LDFLAGS% -o %BINARY-LINUX%

    goto clean

:clean
    set GOOS=
    set GOARCH=

    goto :EOF
