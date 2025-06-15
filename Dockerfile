# Alpine image with WireGuard tools
FROM alpine:latest

RUN apk add --no-cache wireguard-tools bash

COPY entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

ENTRYPOINT ["/entrypoint.sh"]