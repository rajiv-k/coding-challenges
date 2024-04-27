# Build Your Own curl

References:
- https://codingchallenges.fyi/challenges/challenge-curl


## Build & Run

```console
$ make

```

```console
$ build/cccurl -X GET --verbose http://eu.httpbin.org/get
> GET /get HTTP/1.1

< HTTP/1.1 200 200 OK
< Date: Sat, 27 Apr 2024 16:41:21 GMT
< Content-Length: 277
< Connection: keep-alive
< Server: gunicorn/19.9.0
< Access-Control-Allow-Origin: *
< Access-Control-Allow-Credentials: true
< Content-Type: application/json

{
  "args": {},
  "headers": {
    "Accept-Encoding": "gzip",
    "Host": "eu.httpbin.org",
    "User-Agent": "Go-http-client/1.1",
    "X-Amzn-Trace-Id": "Root=1-662d2ab1-1c7f1d0c0537f23b74b66e7d"
  },
  "url": "http://eu.httpbin.org/get"
}

```
