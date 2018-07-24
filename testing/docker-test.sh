#!/bin/bash
# Run tests in golang docker container

go get -d -v
pushd bin
go get -d -v
popd
go test -test.v --running_under_docker
