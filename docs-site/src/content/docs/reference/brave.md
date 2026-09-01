---
title: Brave
description: Brave driven by the chromedriver for its embedded Chromium.
sidebar:
  order: 6
---

| | |
| --- | --- |
| Image | `websummoner/brave` |
| Browser | Brave 1.94.117 (Chromium 152) |
| Driver | ChromeDriver 152.0.7977.64 |

## How it works

Brave's own version line is `major.minor` and unrelated to Chromium's, so the
build tool starts the image, reads the embedded Chromium version from the
browser itself, and fetches the matching chromedriver. Brave 1.94.117 reports
itself as `Brave Browser 152.1.94.117` — Chromium 152.

One detail matters when driving it: `/usr/bin/brave-browser` is a wrapper
script, not the binary. The driver needs the real executable at
`/opt/brave.com/brave/brave-browser`. The hub sets that automatically, so
nothing is needed on the client side.

## Building

```bash
./images brave -t websummoner/brave:1.94
```
