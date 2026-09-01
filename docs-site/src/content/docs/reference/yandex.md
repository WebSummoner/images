---
title: Yandex Browser
description: Yandex Browser with the vendor's own chromedriver fork.
sidebar:
  order: 7
---

| | |
| --- | --- |
| Image | `websummoner/yandex` |
| Browser | Yandex 26.6.1.1083 stable |
| Driver | yandexdriver |

## How it works

`yandexdriver` is Yandex's fork of chromedriver — it even reports itself as
`ChromeDriver` when asked for its version. It is driven exactly like Chrome,
with the browser binary at `/usr/bin/yandex-browser`.

One behaviour is worth knowing when writing tests: shortly after launch Yandex
loads its own start page (`https://ya.ru/`) over whatever the session navigated
to first. A test that navigates immediately can have its page replaced. The
container test suite settles this by visiting `about:blank` at session setup
before the first real navigation.

## Building

```bash
./images yandex -b 26.6.1.1083-1 -t websummoner/yandex:26.6
```
