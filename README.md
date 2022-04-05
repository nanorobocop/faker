# Faker - HTTP server for testing, mocking, faking

## Features

* Handlers for different purposes
* Custom codes, responses
* Prometheus metrics
* Json logs
* Openapi schema validation

## Handlers

* `/echo/EchoString` - repeat EchoString in response
* `/ip` - show client ip
* `/headers` - show client headers
* `/sleep/11` - sleep for some seconds
* `/418` - return HTTP code
* `/Mansur` - will salute you

## Usage

```shell
% ./faker -h
Usage of ./faker:
  -code int
    	response code (default 200)
  -port int
    	port number (default 8080)
  -resp string
    	response content
  -resp-type string
    	response content
  -schema string
    	schema path or url
```

## Example

```shell
$ curl http://localhost:8080/Mansur
Hello there, Mansur!
$ curl http://localhost:8080/418
418 I'm a teapot
$ curl http://localhost:8080/ip
127.0.0.1
$ curl http://localhost:8080/headers
Host: localhost:8080
User-Agent: curl/7.74.0
Accept: */*
```

Logs:

```shell
$ ./faker
Running on port :8080
{"date":"2022-04-05T22:43:59.595680435+09:00","durationNs":135784,"url":"/Mansur","handler":"default","responseCode":200}
{"date":"2022-04-05T22:44:03.641549859+09:00","durationNs":60638,"url":"/418","handler":"code","responseCode":418}
{"date":"2022-04-05T22:44:08.561005332+09:00","durationNs":75179,"url":"/ip","handler":"ip","responseCode":200}
{"date":"2022-04-05T22:44:16.812891108+09:00","durationNs":55359,"url":"/headers","handler":"headers","responseCode":200}
```
