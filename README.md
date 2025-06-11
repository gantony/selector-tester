# Selector Tester

A simple Go web server and web UI for testing selector queries against a Calico Enterprise ui-apis endpoint via a proxy.

## Features

- Web UI (`tool.html` + `tool.js`) for building and sending selector queries.
- Proxy endpoint (`/proxy`) that forwards requests to a target URL, with support for bearer tokens.
- CORS support for browser-based clients.
- Configurable server port via `-port` command-line flag.
- Embedded static assets (no external files needed at runtime).

## Usage

1. **Build and run the server:**
   ```sh
   go run server.go -port 8080
   ```

2. **Open the tool in your browser:**
   ```
   http://127.0.0.1:8080/
   ```

3. **Fill in the form and submit to test your selector queries.**

## Development

- Static files (`tool.html`, `tool.js`) are embedded using Go's `embed` package.
- Need to run server and refresh page to see changes...

## Build

```
./build.sh
```

