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
FROM websummoner/dev_opera:$VERSION

ARG VERSION
LABEL org.websummoner.browser=opera \
      org.websummoner.version=$VERSION

COPY operadriver_linux64/operadriver /usr/bin/
COPY entrypoint.sh /

RUN \
    chmod +x /usr/bin/operadriver && \
    echo "cookie-file = ~/.config/pulse/cookie" >> /etc/pulse/client.conf

USER selenium

EXPOSE 4444
ENTRYPOINT ["/entrypoint.sh"]
