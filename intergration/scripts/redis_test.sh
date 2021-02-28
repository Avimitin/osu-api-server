#!/usr/bin/env sh

export redis_dsn="redis://redis:6379"

echo 'wait for redis set up'
sleep 15

go get ./...
go test -v ./intergration/redis_test.go ./intergration/helper.go
