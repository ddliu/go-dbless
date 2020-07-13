#!/bin/bash
cd test/
env $(grep -v '^#' ../.env | xargs) go test -v