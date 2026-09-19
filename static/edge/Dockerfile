ARG VERSION=latest
FROM golang:1.27-bookworm AS go

COPY devtools /devtools

RUN \
    apt-get update && \
    (apt-get install -y upx-ucl || apt-get install -y upx || echo 'upx unavailable, skipping compression') && \
    cd /devtools && \
    go test -race && \
    GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" && \
    (command -v upx >/dev/null 2>&1 && upx /devtools/devtools || echo 'skipping upx /devtools/devtools')

ARG VERSION=latest
FROM websummoner/dev_edge:$VERSION

ARG VERSION
LABEL org.websummoner.browser=edge \
      org.websummoner.version=$VERSION

ENV DBUS_SESSION_BUS_ADDRESS=/dev/null
COPY --from=go /devtools/devtools /usr/bin/
COPY msedgedriver /usr/bin/
COPY entrypoint.sh /

RUN \
    chmod +x /usr/bin/msedgedriver && \
    echo "cookie-file = ~/.config/pulse/cookie" >> /etc/pulse/client.conf

USER selenium

EXPOSE 4444
ENTRYPOINT ["/entrypoint.sh"]
