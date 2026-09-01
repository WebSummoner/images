#!/bin/bash
SCREEN_RESOLUTION=${SCREEN_RESOLUTION:-"1920x1080x24"}
DISPLAY_NUM=99
export DISPLAY=":$DISPLAY_NUM"

QUIET=${QUIET:-""}
DRIVER_ARGS=""
if [ -z "$QUIET" ]; then
    DRIVER_ARGS="--verbose"
fi

clean() {
  if [ -n "$FILESERVER_PID" ]; then
    kill -TERM "$FILESERVER_PID"
  fi
  if [ -n "$XSELD_PID" ]; then
    kill -TERM "$XSELD_PID"
  fi
  if [ -n "$XVFB_PID" ]; then
    kill -TERM "$XVFB_PID"
  fi
  if [ -n "$DRIVER_PID" ]; then
    kill -TERM "$DRIVER_PID"
  fi
  if [ -n "$PRISM_PID" ]; then
    kill -TERM "$PRISM_PID"
  fi
  if [ -n "$X11VNC_PID" ]; then
    kill -TERM "$X11VNC_PID"
  fi
  if [ -n "$PULSE_PID" ]; then
    kill -TERM "$PULSE_PID"
  fi
}

trap clean SIGINT SIGTERM

/usr/bin/fileserver &
FILESERVER_PID=$!

DISPLAY="$DISPLAY" /usr/bin/xseld &
XSELD_PID=$!

if env | grep -q ROOT_CA_; then
  mkdir -p /tmp/ca-certificates
  for e in $(env | grep ROOT_CA_ | sed -e 's/=.*$//'); do
    certname=$(echo -n $e | sed -e 's/ROOT_CA_//')
    echo ${!e} | base64 -d >/tmp/ca-certificates/${certname}.crt
  done
  update-ca-certificates --localcertsdir /tmp/ca-certificates
fi

/usr/bin/xvfb-run -l -n "$DISPLAY_NUM" -s "-ac -screen 0 $SCREEN_RESOLUTION -noreset -listen tcp" /usr/bin/fluxbox -display "$DISPLAY" >/dev/null 2>&1 &
XVFB_PID=$!

retcode=1
until [ $retcode -eq 0 ]; do
  DISPLAY="$DISPLAY" wmctrl -m >/dev/null 2>&1
  retcode=$?
  if [ $retcode -ne 0 ]; then
    echo Waiting X server...
    sleep 0.1
  fi
done

if [ "$ENABLE_VNC" == "true" ]; then
    x11vnc -display "$DISPLAY" -passwd websummoner -shared -forever -loop500 -rfbport 5900 -rfbportv6 5900 -logfile /dev/null &
    X11VNC_PID=$!
fi

# pulse must share HOME with the browser or in-browser audio never reaches it
mkdir -p ~/.config/pulse
echo -n 'gIvST5iz2S0J1+JlXC1lD3HWvg61vDTV1xbmiGxZnjB6E3psXsjWUVQS4SRrch6rygQgtpw7qmghDFTaekt8qWiCjGvB0LNzQbvhfs1SFYDMakmIXuoqYoWFqTJ+GOXYByxpgCMylMKwpOoANEDePUCj36nwGaJNTNSjL8WBv+Bf3rJXqWnJ/43a0hUhmBBt28Dhiz6Yqowa83Y4iDRNJbxih6rB1vRNDKqRr/J9XJV+dOlM0dI+K6Vf5Ag+2LGZ3rc5sPVqgHgKK0mcNcsn+yCmO+XLQHD1K+QgL8RITs7nNeF1ikYPVgEYnc0CGzHTMvFR7JLgwL2gTXulCdwPbg=='| base64 -d>~/.config/pulse/cookie
pulseaudio --start --exit-idle-time=-1
# wait for the daemon, then expose it over TCP for the video recorder
i=0
until pactl load-module module-native-protocol-tcp >/dev/null 2>&1; do
    i=$((i+1)); [ $i -ge 50 ] && break
    sleep 0.2
done
pactl unload-module module-suspend-on-idle >/dev/null 2>&1 || true
PULSE_PID=$(ps --no-headers -C pulseaudio -o pid | sed -r 's/( )+//g')

DISPLAY="$DISPLAY" /opt/webkit/bin/WebKitWebDriver --port=5555 --host=0.0.0.0 ${DRIVER_ARGS} &
DRIVER_PID=$!

/usr/bin/prism  &
PRISM_PID=$!

wait
