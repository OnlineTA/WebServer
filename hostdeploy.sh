#!/usr/bin/env bash

set -euo pipefail

rsync -av --chown=onlineta:humanta --chmod=D2570,F570 OnlineTA files.json /home/onlineta/

mkdir -p /home/onlineta/www
chown onlineta:humanta /home/onlineta/www
chmod 0774 /home/onlineta/www
rsync -av --chown=onlineta:humanta --chmod=D0774,F464 www /home/onlineta

mkdir -p /home/onlineta/uploads
chown onlineta:humanta /home/onlineta/uploads
chmod 0760 /home/onlineta/uploads

touch /home/onlineta/error.log
chown onlineta:humanta /home/onlineta/error.log
chmod 0644 /home/onlineta/error.log
