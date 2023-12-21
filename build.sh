#!/bin/bash

APP_VERSION=$(git describe --always --long --dirty)

go build -ldflags="-X main.version=${APP_VERSION}" logingestor.go
