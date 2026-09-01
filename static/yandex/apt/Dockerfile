FROM websummoner/browser-base:1.0.0

ARG VERSION=19.4.2.698-1

LABEL browser=yandex-browser-stable:$VERSION

# apt-key is gone since Ubuntu 24.04
RUN  \
        install -d -m 0755 /etc/apt/keyrings && \
        curl -fsSL https://repo.yandex.ru/yandex-browser/YANDEX-BROWSER-KEY.GPG -o /etc/apt/keyrings/yandex-browser.asc && \
        echo 'deb [arch=amd64 signed-by=/etc/apt/keyrings/yandex-browser.asc] https://repo.yandex.ru/yandex-browser/deb stable main' > /etc/apt/sources.list.d/yandex-browser.list && \
        apt-get update && \
        apt-get -y --no-install-recommends install yandex-browser-stable=$VERSION && \
        yandex-browser --version && \
        rm /etc/apt/sources.list.d/yandex-browser.list && \
        rm -Rf /tmp/* && rm -Rf /var/lib/apt/lists/*
