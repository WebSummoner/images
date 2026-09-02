---
title: Firefox
description: Firefox with geckodriver, versioned independently.
sidebar:
  order: 3
---

| | |
| --- | --- |
| Image | `websummoner/firefox` |
| Browser | Mozilla Firefox 155.0.0 |
| Driver | geckodriver 0.37.1 |

## How it works

geckodriver is versioned independently of Firefox and works across a wide range
of releases, so the image takes the newest geckodriver rather than pinning to
the browser build.

The image runs geckodriver directly. Earlier Selenoid-era Firefox images
embedded a copy of the hub binary inside the image; this one does not — that
coupling between two repositories is gone.

geckodriver validates the `Host` header on incoming requests, so the entrypoint
passes `--allow-hosts` and `--allow-origins` for the container's own hostname
and loopback addresses. Without it, requests proxied from the hub are rejected.

## Building

```bash
./images firefox -b 155.0~build1 -t websummoner/firefox:155.0.0
```
