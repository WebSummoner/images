ARG VERSION=latest
FROM websummoner/dev_yandex:$VERSION

ARG VERSION
LABEL org.websummoner.browser=yandex \
      org.websummoner.version=$VERSION

COPY yandexdriver /usr/bin/
COPY entrypoint.sh /

RUN \
    chmod +x /usr/bin/yandexdriver && \
    echo "cookie-file = ~/.config/pulse/cookie" >> /etc/pulse/client.conf

USER selenium

EXPOSE 4444
ENTRYPOINT ["/entrypoint.sh"]
