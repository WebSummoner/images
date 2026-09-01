FROM websummoner/browser-base:1.0.0

ARG VERSION
ARG PACKAGE=google-chrome-stable
ARG INSTALL_DIR=chrome

LABEL browser=$PACKAGE:$VERSION

# apt-key is gone since Ubuntu 24.04
RUN \
        install -d -m 0755 /etc/apt/keyrings && \
        curl -fsSL https://dl-ssl.google.com/linux/linux_signing_key.pub -o /etc/apt/keyrings/google-chrome.asc && \
        echo 'deb [arch=amd64 signed-by=/etc/apt/keyrings/google-chrome.asc] https://dl.google.com/linux/chrome/deb/ stable main' > /etc/apt/sources.list.d/google.list && \
        apt-get update && \
        apt-get -y --no-install-recommends install iproute2 $PACKAGE=$VERSION && \
         ( \
          sed -i -e 's@exec -a "$0" "$HERE/chrome"@& --no-sandbox --disable-gpu@' /opt/google/$INSTALL_DIR/$PACKAGE || \
          sed -i -e 's@exec -a "$0" "$HERE/chrome"@& --no-sandbox --disable-gpu@' /opt/google/$INSTALL_DIR/google-chrome \
         ) && \
        chown root:root /opt/google/$INSTALL_DIR/chrome-sandbox && \
        chmod 4755 /opt/google/$INSTALL_DIR/chrome-sandbox && \
        google-chrome --version && \
        rm -Rf /tmp/* && rm -Rf /var/lib/apt/lists/*
