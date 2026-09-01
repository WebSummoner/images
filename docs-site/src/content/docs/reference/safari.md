---
title: Safari (WebKit)
description: WebKitGTK — the engine Safari is built on — with WebKitWebDriver.
sidebar:
  order: 8
---

| | |
| --- | --- |
| Image | `websummoner/safari` |
| Engine | WebKitGTK 6.0 (`libwebkitgtk-6.0.so.4`) |
| Driver | WebKitWebDriver |

:::caution
Real Safari runs only on macOS and iOS. This image uses **WebKitGTK**, the same
engine Safari is built on. Functionally the two are equivalent, but fonts and
pixel-for-pixel rendering can differ.
:::

## How it works

WebKitWebDriver is built from the same source tree as the engine, so driver and
browser always match — there is no version-pinning problem here.

`WebKitWebDriver` is less complete than the Chromium and Gecko drivers, and the
hub fills several gaps from outside the browser: file upload, and translating
the `proxy` capability into the system environment WebKit actually reads. Two
behaviours still show through to your tests — cookies must set `sameSite`, and
`quit()` can throw even though the session ended. All four are documented on the
hub's
[WebKit section](https://websummoner.github.io/websummoner/reference/browser-images/#webkit-safari).

## Building

```bash
./images safari -t websummoner/safari:2.52.6
```

:::caution
WebKitGTK is compiled from source: expect one to two hours and a large memory
footprint. Build it deliberately, on a machine chosen for it — not as part of a
routine release run.
:::
