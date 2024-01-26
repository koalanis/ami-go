#!/bin/bash
git checkout main
git pull origin main

sudo rm -rf /opt/ami-go
go build
sudo mv ./ami-go /opt/
