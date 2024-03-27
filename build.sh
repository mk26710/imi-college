#!/bin/bash
go build -o college -ldflags "-s -w" cmd/imi/college/main.go
