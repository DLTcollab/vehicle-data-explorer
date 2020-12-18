#!/usr/bin/env bash
git submodule update --init --recursive
flatc --go third_party/vehicle-scheme/*.fbs
go build -v ./...
