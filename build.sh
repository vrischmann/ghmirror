#!/bin/bash

export GO15VENDOREXPERIMENT=1

go build --ldflags="-X main.commit=$(git rev-parse HEAD) -X main.version=$(cat VERSION)"
