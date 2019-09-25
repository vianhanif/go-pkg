#!/bin/bash

go clean -testcache
go test -v -cover -coverprofile=/tmp/coverage.out -p 1 ./...