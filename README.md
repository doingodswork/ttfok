# ttfok

`ttfok`: **T**ime **t**o **f**irst **OK**

CLI to measure the startup time of a web service until it serves a request

## Usage

`ttfok /path/to/your/app http://localhost:8080/some/path`

Or with app arguments:

`ttfok /path/to/your/app --some arg --another arg http://localhost:8080/some/path`

To customize times:

```text
  -t duration
        Timeout for the request (default 1ms)
  -w duration
        Duration to wait for app start (default 1s)
```

### Example

`ttfok -w 2s node ./addon.js http://localhost:7000/manifest.json`
