#!/usr/bin/env sh

chown -R retriever /init-out || echo "Could not chown /init-out"

if [ "$(id -u)" = '0' ]; then
    su-exec retriever "$@"
fi 