#!/bin/bash
git pull
go get
go build main.go
mv main ~/.local/bin/ror
