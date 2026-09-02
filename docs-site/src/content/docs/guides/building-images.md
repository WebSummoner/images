---
title: Building images
description: How the layers fit together and how to build each browser with the images tool.
sidebar:
  order: 1
---

## What's inside an image

Each image consists of three layers:

1. **Base layer** — the things every image needs: Xvfb, fonts, locales, the
   cursor blinking fix, timezone definition and so on. Built once and shared.
2. **Browser layer** — the browser binary.
3. **Driver layer** — the matching WebDriver binary.

## Building procedure

The procedure is automated by a Go binary with all Docker build files embedded.
To build it from source:

```bash
go build
```

To show help:

```bash
./images --help
./images firefox --help
```

Before building you can optionally clone the tests repository next to this one:

```bash
git clone https://github.com/WebSummoner/websummoner-container-tests.git
```

```
images/                        # this repo
websummoner-container-tests/   # optional tests repo
```

Then add `--test` to run the suite against the image you just built.

### Firefox

```bash
./images firefox -b 155.0~build1 -t websummoner/firefox:155.0.0
```

Omit `-d` (or pass `latest`) to take the newest matching geckodriver.

### Chrome

```bash
./images chrome -b 152.0.7977.75-1 -d latest -t websummoner/chrome:152.0
```

### Microsoft Edge

```bash
./images edge -b 152.0.4191.62-1 -d 152.0.4191.62 -t websummoner/edge:152.0
```

### Opera

```bash
./images opera -b 135.0.5973.76 -t websummoner/opera:135.0
```

Do not pass `-d` for Opera. `operachromiumdriver` tags follow the Chromium
version Opera is built on, not Opera's own line, so the tool resolves it — see
[Opera](https://websummoner.riadvice.com/websummoner/reference/browser-images/#opera).

### Yandex Browser

```bash
./images yandex -b 26.6.1.1083-1 -t websummoner/yandex:26.6
```

### Brave

```bash
./images brave -t websummoner/brave:1.94
```

Brave's line is major.minor; the tool detects the embedded Chromium from the
built image and fetches the matching chromedriver.

### WebKit (Safari)

```bash
./images safari -t websummoner/safari:2.52.6
```

:::caution
WebKitGTK is compiled from source and takes one to two hours with a large
memory footprint. Build it deliberately, on a machine chosen for it.
:::

## Building from a local package

Replace the package version with a path to a `.deb`:

```bash
./images firefox -b /path/to/firefox_155.0_amd64.deb -t websummoner/firefox:155.0.0
```

The file name must contain the full version — the tool derives the browser
version by parsing it.

## Pushing

Add `--push` to push after building:

```bash
./images chrome -b 152.0.7977.75-1 -t websummoner/chrome:152.0 --push
```

## Selecting a browser channel

Besides the default stable channel:

| Browser | Channel | Package |
| --- | --- | --- |
| firefox | beta | `firefox` ([PPA](http://launchpad.net/~mozillateam/+archive/firefox-next/+packages)) |
| firefox | dev | `firefox-trunk` ([PPA](http://launchpad.net/~ubuntu-mozilla-daily/+archive/ppa/+packages)) |
| firefox | esr | `firefox-esr` ([PPA](http://launchpad.net/~mozillateam/+archive/ppa/+packages)) |
| chrome | beta | `google-chrome-beta` |
| chrome | dev | `google-chrome-unstable` |
| opera | beta | `opera-beta` |
| opera | dev | `opera-developer` |

```bash
./images firefox -b <version> --channel dev -t websummoner/firefox:dev
```
