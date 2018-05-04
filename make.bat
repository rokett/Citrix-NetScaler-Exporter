@echo off
SETLOCAL

set _TARGETS=build

set APP=Citrix-NetScaler-Exporter
set VERSION=3.1.0
set BINARY-X86=%APP%_%VERSION%_Windows_32bit.exe
set BINARY-X64=%APP%_%VERSION%_Windows_64bit.exe

REM Set build number from git commit hash
for /f %%i in ('git rev-parse HEAD') do set BUILD=%%i

if [%1]==[] goto usage

REM *** CHECK THAT VALID ARG IS PASSED ***

set LDFLAGS=-ldflags "-X main.version=%VERSION% -X main.build=%BUILD%"

goto %1

:build
    set GOOS=windows

    echo "=== Building x86 ==="
    set GOARCH=386

    go build -o %BINARY-X86% %LDFLAGS%

    echo "=== Building x64 ==="
    set GOARCH=amd64

    go build -o %BINARY-X64% %LDFLAGS%

    goto :finalise

:usage
    echo usage: make [target]
    echo.
    echo target is one of {%_TARGETS%}.
    exit /b 2
    goto :eof

:finalise
    set GOOS=
    set GOARCH=

    goto :EOF
