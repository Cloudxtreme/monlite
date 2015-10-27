#! /bin/bash

GOOS=linux GOARCH=amd64 gb build -r all || exit 1

mv bin/mon-linux-amd64 linux/monlite.linux

rm -f linux/*.xz

xz -z9 linux/monlite.linux

exit 0
