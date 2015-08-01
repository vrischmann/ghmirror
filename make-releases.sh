#!/bin/bash

export GO15VENDOREXPERIMENT=1

gox --os="linux windows darwing openbsd freebsd" \
    --arch="386 amd64" \
    --output="build/{{.Dir}}_{{.OS}}_{{.Arch}}"
