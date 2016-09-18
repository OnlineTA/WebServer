#!/usr/bin/env bash

set -euo pipefail

. deploy.conf

rsync -av OnlineTA files.json www hostdeploy.sh ${USER}@${HOST}:~${USER}/
