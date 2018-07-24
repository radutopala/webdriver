#!/bin/bash
# Run tests under Travis for continuous integration.

which chrome
which google-chrome-stable
which google-chrome

go test -coverprofile=coverage.txt -covermode=atomic -test.v -timeout=20m -chrome_binary=google-chrome-stable
