#!/bin/sh

set -e

ssh -i "${HETZNER_ACCESS_SSH}" kaspars@95.216.217.197 docker-compose up -d --no-deps --build web

ssh -i "${HETZNER_ACCESS_SSH}" kaspars@95.216.217.197 docker-compose up -d --no-deps --build api
