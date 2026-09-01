ARG UBUNTU_VERSION=26.04

FROM ubuntu:$UBUNTU_VERSION AS build

ARG WEBKIT_VERSION="2.52.6"

RUN \
    apt-get update && \
    apt-get -y install --no-install-recommends ca-certificates curl xz-utils && \
    curl --retry 8 --retry-all-errors --retry-delay 5 -fSLo /webkitgtk.tar.xz https://webkitgtk.org/releases/webkitgtk-$WEBKIT_VERSION.tar.xz && \
    mkdir -p /webkit /opt/webkit && \
    tar -xJf /webkitgtk.tar.xz -C /webkit --strip-components=1 && \
    mkdir -p /webkit/Tools/glib/dependencies && \
    curl --retry 8 --retry-all-errors --retry-delay 5 -fSLo /webkit/Tools/glib/dependencies/apt https://raw.githubusercontent.com/WebKit/WebKit/wpewebkit-$WEBKIT_VERSION/Tools/glib/dependencies/apt && \
    sed -i 's|    g_return_val_if_fail(WEBKIT_IS_WEB_VIEW(view), NULL);|    g_return_val_if_fail(WEBKIT_IS_WEB_VIEW(view), NULL);\n    webkit_settings_set_media_playback_requires_user_gesture(webkit_web_view_get_settings(view), FALSE);|' /webkit/Tools/MiniBrowser/gtk/BrowserTab.c && \
    cd /webkit && \
    yes | DEBIAN_FRONTEND=noninteractive Tools/gtk/install-dependencies && \
    cmake -DPORT=GTK -DCMAKE_BUILD_TYPE=Release -DCMAKE_INSTALL_PREFIX=/opt/webkit -DUSE_WPE_RENDERER=OFF -DENABLE_MINIBROWSER=ON -DENABLE_BUBBLEWRAP_SANDBOX=OFF -DENABLE_GAMEPAD=OFF -DENABLE_SPELLCHECK=OFF -DENABLE_WAYLAND_TARGET=OFF -DUSE_OPENJPEG=OFF -DUSE_LIBBACKTRACE=OFF -GNinja && \
    ninja && \
    ninja install && \
    rm -Rf /var/lib/apt/lists/*

FROM golang:1.27-bookworm AS go

COPY cmd/prism /prism

RUN \
    apt-get update && \
    apt-get install -y libx11-dev && (apt-get install -y upx-ucl || apt-get install -y upx || echo 'upx unavailable, skipping compression') && \
    cd /prism && \
    GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" && \
    (command -v upx >/dev/null 2>&1 && upx /prism/prism || echo 'skipping upx /prism/prism') && \
    rm -Rf /var/lib/apt/lists/*

FROM websummoner/browser-base:1.0.0

ARG WEBKIT_VERSION

LABEL browser=safari:$WEBKIT_VERSION

ENV WEBKIT_VERSION=$WEBKIT_VERSION
ENV LD_LIBRARY_PATH=/opt/webkit/lib/

COPY --from=build /opt/webkit /opt/webkit
COPY --from=go /prism/prism /usr/bin/

RUN \
    apt-get update && \
    apt-get -y install --no-install-recommends \
        libgtk-3-0t64 \
        libgtk-4-1 \
        libavif16 \
        libsoup-3.0-0 \
        libwebp7 \
        libwebpdemux2 \
        libsecret-1-0 \
        libhyphen0 \
        libwoff1 \
        libharfbuzz-icu0 \
        libgstreamer-gl1.0-0 \
        gstreamer1.0-plugins-good \
        gstreamer1.0-libav \
        libopenjp2-7 \
        libnotify4 \
        libxslt1.1 \
        libegl1 \
        libatomic1 \
        libopengl0 && \
    ldconfig && \
    apt-get clean && \
    rm -Rf /tmp/* && rm -Rf /var/lib/apt/lists/* && \
    chmod 777 /etc/ssl/certs && \
    mkdir /tmp/ca-certificates && \
    chmod 777 /tmp/ca-certificates

COPY entrypoint.sh /

USER selenium

EXPOSE 4444
ENTRYPOINT ["/entrypoint.sh"]
