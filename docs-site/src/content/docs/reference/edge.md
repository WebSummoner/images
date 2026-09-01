---
title: Microsoft Edge
description: Edge with an exactly-matched Edge WebDriver.
sidebar:
  order: 4
---

| | |
| --- | --- |
| Image | `websummoner/edge` |
| Browser | Microsoft Edge 152.0.4191.53 |
| Driver | Microsoft Edge WebDriver 152.0.4191.53 |

## How it works

Edge is Chromium, and msedgedriver behaves like chromedriver: it enforces an
exact version match. Microsoft's CDN serves a driver for each exact browser
build, so the image pins both to the same version.

## Building

```bash
./images edge -b 152.0.4191.53-1 -d 152.0.4191.53 -t websummoner/edge:152.0
```
