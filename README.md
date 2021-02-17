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

Later, open the following url: http://localhost:8080
