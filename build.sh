#!/usr/bin/env bash

set -euo pipefail

go build --ldflags '-linkmode external -extldflags "-static"'
