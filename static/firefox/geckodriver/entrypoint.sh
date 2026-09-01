#!/bin/bash
SCREEN_RESOLUTION=${SCREEN_RESOLUTION:-"1920x1080x24"}
DISPLAY_NUM=99
export DISPLAY=":$DISPLAY_NUM"

VERBOSE=${VERBOSE:-""}
DRIVER_ARGS=${DRIVER_ARGS:-""}
if [ -n "$VERBOSE" ]; then
    DRIVER_ARGS="$DRIVER_ARGS, \"--log\", \"debug\""
fi

clean() {
  if [ -n "$FILESERVER_PID" ]; then
    kill -TERM "$FILESERVER_PID"
  fi
  if [ -n "$XSELD_PID" ]; then
    kill -TERM "$XSELD_PID"
  fi
  if [ -n "$PULSE_PID" ]; then
    kill -TERM "$PULSE_PID"
  fi
  if [ -n "$XVFB_PID" ]; then
    kill -TERM "$XVFB_PID"
  fi
  if [ -n "$DRIVER_PID" ]; then
    kill -TERM "$DRIVER_PID"
  fi
  if [ -n "$X11VNC_PID" ]; then
    kill -TERM "$X11VNC_PID"
  fi
}

trap clean SIGINT SIGTERM

/usr/bin/fileserver &
FILESERVER_PID=$!

DISPLAY="$DISPLAY" /usr/bin/xseld &
XSELD_PID=$!

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

/usr/bin/xvfb-run -l -n "$DISPLAY_NUM" -s "-ac -screen 0 $SCREEN_RESOLUTION -noreset -listen tcp" /usr/bin/fluxbox -display "$DISPLAY" -log /dev/null 2>/dev/null &
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

# geckodriver validates the Host header, and the container hostname is only
# known at run time, so it has to be allowed explicitly.
DISPLAY="$DISPLAY" /usr/bin/geckodriver --host 0.0.0.0 --port 4444 \
  --allow-hosts "$(hostname)" localhost 127.0.0.1 "$(hostname -i 2>/dev/null || echo 127.0.0.1)" \
  --allow-origins "http://$(hostname):4444" "http://localhost:4444" ${DRIVER_ARGS} &
DRIVER_PID=$!

if env | grep -q ROOT_CA_; then
  while true; do
    if certDB=$(ls -d /tmp/rust_mozprofile*/cert9.db 2>/dev/null); then
      break
    else
      sleep 0.1
    fi
  done
  certdir=$(dirname ${certDB})
  for e in $(env | grep ROOT_CA_ | sed -e 's/=.*$//'); do
    certname=$(echo -n $e | sed -e 's/ROOT_CA_//')
    echo ${!e} | base64 -d >/tmp/cert.pem
    certutil -A -n ${certname} -t "TC,C,T" -i /tmp/cert.pem -d sql:${certdir}
    if cat tmp/cert.pem | grep -q "PRIVATE KEY"; then
      PRIVATE_KEY_PASS=${PRIVATE_KEY_PASS:-\'\'}
      openssl pkcs12 -export -in /tmp/cert.pem -clcerts -nodes -out /tmp/key.p12 -passout pass:${PRIVATE_KEY_PASS}  -passin pass:${PRIVATE_KEY_PASS}
      pk12util -d sql:${certdir} -i /tmp/key.p12 -W ${PRIVATE_KEY_PASS}
      rm /tmp/key.p12
    fi
    rm /tmp/cert.pem
  done
fi

wait
