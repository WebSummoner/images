ARG VERSION=latest
FROM websummoner/dev_firefox:$VERSION

COPY geckodriver /usr/bin/
COPY entrypoint.sh /

RUN echo "cookie-file = ~/.config/pulse/cookie" >> /etc/pulse/client.conf

USER selenium

EXPOSE 4444
ENTRYPOINT ["/entrypoint.sh"]
