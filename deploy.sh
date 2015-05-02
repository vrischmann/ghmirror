#!/bin/bash

gox --osarch=linux/amd64

rsync -avz --delete . sphax@server.rischmann.fr:/data/ghmirror/.

ssh root@server.rischmann.fr "systemctl restart ghmirror"
