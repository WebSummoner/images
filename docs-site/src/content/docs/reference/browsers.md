---
title: Browsers
description: The seven browsers built here, the drivers they ship with, and how each is driven.
sidebar:
  order: 1
---

Every image is built on Ubuntu 26.04 LTS from a shared base layer, then a
browser layer, then a driver layer.

| | Browser | Version | Driver | Driver matches |
| --- | --- | --- | --- | --- |
| <img src="https://raw.githubusercontent.com/alrra/browser-logos/main/src/chrome/chrome.svg" width="20" alt="Chrome"> | [Chrome](/reference/chrome/) | 152.0.7977.64 | ChromeDriver 152.0.7977.64 | exact build |
| <img src="https://raw.githubusercontent.com/alrra/browser-logos/main/src/firefox/firefox.svg" width="20" alt="Firefox"> | [Firefox](/reference/firefox/) | 155.0.0 | geckodriver 0.37.1 | independent |
| <img src="https://raw.githubusercontent.com/alrra/browser-logos/main/src/edge/edge.svg" width="20" alt="Edge"> | [Edge](/reference/edge/) | 152.0.4191.53 | Edge WebDriver 152.0.4191.53 | exact build |
| <img src="https://raw.githubusercontent.com/alrra/browser-logos/main/src/opera/opera.svg" width="20" alt="Opera"> | [Opera](/reference/opera/) | 135.0.5973.66 | OperaDriver 150.0.7871.212 | Chromium line |
| <img src="https://raw.githubusercontent.com/alrra/browser-logos/main/src/brave/brave.svg" width="20" alt="Brave"> | [Brave](/reference/brave/) | 1.94.117 (Chromium 152) | ChromeDriver 152.0.7977.64 | embedded Chromium |
| <img src="https://raw.githubusercontent.com/alrra/browser-logos/main/src/yandex/yandex_24x24.png" width="20" alt="Yandex"> | [Yandex](/reference/yandex/) | 26.6.1.1083 | yandexdriver | vendor fork |
| <img src="https://raw.githubusercontent.com/alrra/browser-logos/main/src/safari/safari.svg" width="20" alt="Safari"> | [Safari](/reference/safari/) | WebKitGTK 6.0 | WebKitWebDriver | ships together |

The versions above are what the current images contain. Tags and the versioning
policy are documented in
[Image tags](https://websummoner.riadvice.com/websummoner/reference/image-tags/).

## How a browser is driven

Every image runs a WebDriver server on port 4444 and the browser beside it. The
differences are in how tightly the driver is coupled to the browser build:

- **Chromium-family browsers** (Chrome, Edge, Brave, Yandex, Opera) speak the
  same protocol through a chromedriver derivative. The driver enforces a
  version check against the browser, which is why each image pins a specific
  driver rather than tracking latest.
- **Firefox** uses geckodriver, versioned independently of Firefox and
  compatible across a wide range of releases.
- **WebKit** uses WebKitWebDriver, built from the same source tree as the
  engine, so the two always match.
