FROM websummoner/browser-base:1.0.0

ARG VERSION
ARG PACKAGE=firefox

LABEL browser=$PACKAGE:$VERSION

# Mozilla's official APT repository (the old ubuntu-mozillateam PPAs are EOL):
# https://support.mozilla.org/en-US/kb/install-firefox-linux#w_install-firefox-deb-package-for-debian-based-distributions-recommended
RUN  \
        apt-get update && \
        apt-get install -y software-properties-common wget && \
        install -d -m 0755 /etc/apt/keyrings && \
        wget -q https://packages.mozilla.org/apt/repo-signing-key.gpg -O- | tee /etc/apt/keyrings/packages.mozilla.org.asc > /dev/null && \
        echo "deb [signed-by=/etc/apt/keyrings/packages.mozilla.org.asc] https://packages.mozilla.org/apt mozilla main" | tee /etc/apt/sources.list.d/mozilla.list > /dev/null && \
        printf 'Package: *\nPin: origin packages.mozilla.org\nPin-Priority: 1000\n' | tee /etc/apt/preferences.d/mozilla > /dev/null && \
        apt-get update && \
        apt-get -y --no-install-recommends install iproute2 "$(apt-cache search --names-only '^libavcodec[0-9]+$' | awk '{print $1}' | sort -V | tail -1)" $PACKAGE=$VERSION && \
        ( [ "$PACKAGE" != "firefox" ] && ln /usr/bin/$PACKAGE /usr/bin/firefox ) || true && \
        firefox --version && \
        rm -Rf /tmp/* && rm -Rf /var/lib/apt/lists/*
