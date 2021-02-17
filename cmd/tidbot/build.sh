#!/bin/bash
VERSION=$(git describe --tags)
echo Building version $VERSION
go build -ldflags "-X github.com/Pedro-Pessoa/tidbot/common.VERSION=${VERSION}"