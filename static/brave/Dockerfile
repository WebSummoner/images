ARG VERSION=latest
FROM golang:1.27-bookworm AS go

COPY devtools /devtools

RUN \
    apt-get update && \
    apt-get install -y libx11-dev && (apt-get install -y upx-ucl || apt-get install -y upx || echo 'upx unavailable, skipping compression') && \
    cd /devtools && \
    go test -race && \
    GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" && \
    (command -v upx >/dev/null 2>&1 && upx /devtools/devtools || echo 'skipping upx /devtools/devtools')

ARG VERSION=latest
FROM websummoner/dev_brave:$VERSION

ENV DBUS_SESSION_BUS_ADDRESS=/dev/null
COPY --from=go /devtools/devtools /usr/bin/
COPY policies.json /etc/brave/policies/managed/
COPY chromedriver /usr/bin/
COPY entrypoint.sh /

# chromedriver looks for google-chrome on PATH by default — give it a
# launcher that runs the real Brave binary with the flags a root container
# needs (the distro brave-browser wrapper mis-resolves its location here).
RUN printf '#!/bin/sh\nexec /opt/brave.com/brave/brave-browser --no-sandbox --disable-gpu "$@"\n' > /usr/bin/google-chrome && chmod 755 /usr/bin/google-chrome && \
    chmod +x /usr/bin/chromedriver && \
    chown -R root:root /etc/brave/policies/managed/ && \
    echo "cookie-file = ~/.config/pulse/cookie" >> /etc/pulse/client.conf

USER root

EXPOSE 4444
ENTRYPOINT ["/entrypoint.sh"]
