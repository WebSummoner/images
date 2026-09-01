---
title: Quick start
description: Pull a published image, or build one yourself with the images tool.
---

## Using a published image

Most people never build anything — pull the image and point
[WebSummoner](https://websummoner.github.io/websummoner/) at it:

```bash
docker pull websummoner/chrome:latest
```

Available tags and the versioning policy are documented in
[Image tags](https://websummoner.github.io/websummoner/reference/image-tags/).

## Building the tool

The build tool is a single Go binary; all Docker build files are embedded in it.

```bash
go build
./images --help
./images firefox --help
```

Then see [Building images](/guides/building-images/).
