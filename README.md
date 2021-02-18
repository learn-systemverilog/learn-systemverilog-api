# Learn SystemVerilog API
[![Go Report Card](https://goreportcard.com/badge/github.com/learn-systemverilog/learn-systemverilog-api)](https://goreportcard.com/report/github.com/learn-systemverilog/learn-systemverilog-api)
[![DeepSource](https://deepsource.io/gh/learn-systemverilog/learn-systemverilog-api.svg/?label=active+issues&token=mjKw9zrb9k0KlMHUmAHtlFIe)](https://deepsource.io/gh/learn-systemverilog/learn-systemverilog-api/?ref=repository-badge)

Learn SystemVerilog API is the API used by [learn-systemverilog-web](https://github.com/learn-systemverilog/learn-systemverilog-web). Currently, it transpiles the code written in SystemVerilog to JavaScript so that the simulation can work in any browser.

## Requirements
- [Docker Desktop](https://www.docker.com/products/docker-desktop)

## Setup
First of all, clone the repository:
```bash
git clone https://github.com/learn-systemverilog/learn-systemverilog-api.git

cd learn-systemverilog-api
```

Then, build the docker image:
```bash
docker build --tag learn-systemverilog-api .
```

Next, run the docker image as a container:
```bash
docker run --publish 8080:8080 learn-systemverilog-api
```

Later, open this url: http://localhost:8080. You should see the following message:
```
{"message":"Hello, World!"}
```
Congratulations!

## Usage
To transpile a code written in SystemVerilog to JavaScript, you will need to watch (you can use your browser) for the Server-sent events on the following endpoint:
```
GET http://localhost:8080/transpile?code={YOUR_SYSTEMVERILOG_CODE}
```

There are four types of events:
- `internal`: Internal logs from the server.
  - Format: `{"message": "...", "severity": "debug|info|warn|error"}`
- `stdout`: Standard output written by the transpilers.
  - Format: `{"stdout": "..."}`
- `stderr`: Standard error written by the transpilers.
  - Format: `{"stderr": "..."}`
- `output`: A JSON encoded string representing the transpiled JavaScript code.
  - Format: `var Module = typeof Module !== 'undefined' ? Module : {};\n\n// --...`
