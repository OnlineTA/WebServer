#!/usr/bin/env bash

set -euo pipefail

go build -o OnlineTA --ldflags '-linkmode external -extldflags "-static"'
