---
title: Chrome
description: Google Chrome with an exactly-matched ChromeDriver.
sidebar:
  order: 2
---

| | |
| --- | --- |
| Image | `websummoner/chrome` |
| Browser | Google Chrome 152.0.7977.75 |
| Driver | ChromeDriver 152.0.7977.75 |

## How it works

ChromeDriver enforces an exact major-version match against the browser, so the
image pins the driver to the same build as Chrome. Both come from
Chrome for Testing, which publishes browser and driver together for every
release — so the pair is always consistent.

The driver launches Chrome directly; nothing in the hub rewrites the browser
name or binary path for it.

## Building

```bash
./images chrome -b 152.0.7977.75-1 -d latest -t websummoner/chrome:152.0
```

`-d latest` takes the matching Chrome for Testing driver. Pass an explicit
version only when reproducing an older image.
