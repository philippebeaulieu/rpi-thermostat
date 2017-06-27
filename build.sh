#!/bin/bash

export GOPATH=/usr/local/rpi-thermostat
cd /usr/local/rpi-thermostat/src/github.com/philippebeaulieu/rpi-thermostat
git pull
GOOS=linux GOARCH=arm go build
./rpi-thermostat