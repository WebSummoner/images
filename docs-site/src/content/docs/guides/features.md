---
title: Image features
description: Custom root certificates and custom Chrome profiles.
sidebar:
  order: 2
---

## Adding a custom root certification authority

In corporate networks, tested environments often use self-signed
[TLS](https://en.wikipedia.org/wiki/Transport_Layer_Security) certificates
issued by a [root certification authority](https://en.wikipedia.org/wiki/Root_certificate)
that browsers do not know. A browser refuses such a page with "Your connection
is not private". The `acceptInsecureCerts` capability ignores certificate
errors, but it does not help when the page uses
[HSTS](https://en.wikipedia.org/wiki/HTTP_Strict_Transport_Security).

To work properly, add your root certificate to the trusted list. These images
accept it through an environment variable. For a certificate in `cert.pem`:

1. Encode it as [Base64](https://en.wikipedia.org/wiki/Base64):

   ```bash
   CERT_CONTENTS=$(cat cert.pem | base64 -w0)
   ```

   On macOS:

   ```bash
   CERT_CONTENTS=$(cat cert.pem | base64)
   ```

2. Pass it to the browser image:

   ```
   ROOT_CA_<cert-name>="$CERT_CONTENTS"
   ```

   `<cert-name>` becomes the certificate name in the browser's certificate
   store, for example `ROOT_CA_MY_CERT="LS0tL....=="`.

## Using a custom browser profile with Chrome

When launching Chrome with a custom profile directory, DevTools will not work
unless you also set `BROWSER_PROFILE_DIR`:

```json
{
  "capabilities": {
    "alwaysMatch": {
      "browserName": "chrome",
      "browserVersion": "152.0",
      "goog:chromeOptions": {
        "args": [ "user-data-dir=/profiles/custom.XYZ" ]
      },
      "websummoner:options": {
        "env": [ "BROWSER_PROFILE_DIR=/profiles/custom.XYZ" ]
      }
    }
  }
}
```
