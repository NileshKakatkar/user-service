#!/bin/sh
set -e
if [ -n "${JWT_PRIVATE_KEY_PATH:-}" ] && [ -r "$JWT_PRIVATE_KEY_PATH" ]; then
    cp "$JWT_PRIVATE_KEY_PATH" /tmp/shopcore-private.pem
    chown 10001:10001 /tmp/shopcore-private.pem
    chmod 600 /tmp/shopcore-private.pem
    export JWT_PRIVATE_KEY_PATH=/tmp/shopcore-private.pem
fi
exec su-exec 10001:10001 /app/service
