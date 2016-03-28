#!/bin/bash

gox --ldflags="-X main.commit=$(git rev-parse HEAD) -X main.version=$(cat VERSION)" \
    --os="linux windows darwing openbsd freebsd" \
    --arch="386 amd64" \
    --output="build/{{.Dir}}_{{.OS}}_{{.Arch}}"
