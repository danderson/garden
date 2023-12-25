#!/bin/bash

set -eu

mkdir -p /state/tailscale /state/garden
chmod 700 /state/tailscale /state/garden

./bin/migrate
exec ./proxy
