#!/usr/bin/env bash

set -euo pipefail

. deploy.conf

rsync -av OnlineTA www onlineta@${HOST}:~onlineta/
