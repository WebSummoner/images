---
title: Opera
description: Opera with its own driver, resolved from the Chromium line.
sidebar:
  order: 5
---

| | |
| --- | --- |
| Image | `websummoner/opera` |
| Browser | Opera 135.0.5973.76 (Chromium 151) |
| Driver | OperaDriver 150.0.7871.212 |

## How it works

Opera is the awkward one, for two reasons.

**Driver tags follow Chromium, not Opera.** `operachromiumdriver` releases are
numbered by the Chromium version Opera is built on — Opera N ships Chromium
N+16. Passing `-d 135.0.…` fetches a driver sixteen majors too old, which
refuses every session. The tool works out the Chromium line itself, so do not
pass `-d`.

**Opera publishes late.** At the time of writing there is no driver on the
Chromium 151 line that Opera 135 needs, so the tool falls back to the newest
`operadriver` — never a chromedriver. The version check in this driver family
is a *warning*, not a refusal: OperaDriver 150 drives Opera 135 and logs

```
This version of OperaDriver has not been tested with Opera version 151.
```

then works normally. Substituting a Chrome for Testing chromedriver also starts
sessions, but crashes the renderer whenever a page opens a window. Opera patches
its Chromium and only its own driver accounts for that — a real driver one line
behind beats a foreign driver on the exact line.

Two more behaviours show through to tests, documented on the hub's
[Opera section](https://websummoner.riadvice.com/websummoner/reference/browser-images/#opera):
operadriver answers in legacy JSONWP unless asked for W3C (the hub asks), and a
fresh session reports Opera's own UI pages as extra window handles.

## Building

```bash
./images opera -b 135.0.5973.76 -t websummoner/opera:135.0
```
