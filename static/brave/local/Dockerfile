FROM websummoner/browser-base:1.0.0

ARG VERSION
ARG PACKAGE=brave-browser

LABEL browser=$PACKAGE:$VERSION

RUN \
        curl -fsSLo brave-browser.deb http://host.docker.internal:8080/brave-browser.deb && \
        apt-get update && \
        apt-get -y --no-install-recommends install iproute2 ./brave-browser.deb && \
        ( \
          sed -i -e 's@exec -a "$0" "$HERE/brave"@& --no-sandbox --disable-gpu@' /usr/bin/brave-browser || true \
        ) && \
        mkdir -p /etc/brave/policies/managed && \
        chown -R root:root /etc/brave/policies && \
        brave-browser --version && \
        rm -Rf /tmp/* && rm -Rf /var/lib/apt/lists/* && rm -f brave-browser.deb
