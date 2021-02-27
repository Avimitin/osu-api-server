#!/usr/bin/env sh

export database_host="database:3306"
export osuapi_project_root="/go/src/github.com/avimitin/osuapi"

echo "wait for database setting up"
sleep 15

go get ./...
go test -v ./intergration/mysql_test.go ./intergration/helper.go
