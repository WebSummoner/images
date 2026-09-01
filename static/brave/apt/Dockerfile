FROM websummoner/browser-base:1.0.0

ARG VERSION
ARG PACKAGE=brave-browser

LABEL browser=$PACKAGE:$VERSION

RUN \
        curl -fsSLo /usr/share/keyrings/brave-browser-archive-keyring.gpg https://brave-browser-apt-release.s3.brave.com/brave-browser-archive-keyring.gpg && \
        echo 'deb [arch=amd64 signed-by=/usr/share/keyrings/brave-browser-archive-keyring.gpg] https://brave-browser-apt-release.s3.brave.com/ stable main' > /etc/apt/sources.list.d/brave-browser-release.list && \
        apt-get update && \
        apt-get -y --no-install-recommends install iproute2 $PACKAGE=$VERSION && \
         ( \
          sed -i -e 's@exec -a "$0" "$HERE/brave"@& --no-sandbox --disable-gpu@' /usr/bin/brave-browser || true \
         ) && \
        mkdir -p /etc/brave/policies/managed && \
        chown -R root:root /etc/brave/policies && \
        brave-browser --version && \
        rm -Rf /tmp/* && rm -Rf /var/lib/apt/lists/*
