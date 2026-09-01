FROM websummoner/browser-base:1.0.0

ARG VERSION
ARG PACKAGE=opera-stable

LABEL browser=$PACKAGE:$VERSION

# apt-key is gone since Ubuntu 24.04
RUN  \
        install -d -m 0755 /etc/apt/keyrings && \
        curl -fsSL https://deb.opera.com/archive.key -o /etc/apt/keyrings/opera.asc && \
        echo 'deb [signed-by=/etc/apt/keyrings/opera.asc] https://deb.opera.com/opera-stable/ stable non-free' > /etc/apt/sources.list.d/opera-blink.list && \
        apt-get update && \
        apt-get -y --no-install-recommends install libgtk-3-0 $PACKAGE=$VERSION && \
        ( [ "$PACKAGE" != "opera-stable" ] && ln /usr/bin/$PACKAGE /usr/bin/opera ) || true && \
        opera --version && \
        rm -f /etc/apt/sources.list.d/$PACKAGE.list /etc/apt/sources.list.d/opera-blink.list && \
        rm -Rf /tmp/* && rm -Rf /var/lib/apt/lists/*
