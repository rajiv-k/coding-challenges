# Build Your Own Netcat

References:
- https://codingchallenges.fyi/challenges/challenge-netcat

## Build & Run

```console
$ make
$ build/ccnc -l -p 6969 -e /bin/bash
```

Then connect to this server from a different terminal session.

```console
$ echo date | nc localhost 6969
Sat Apr 27 17:12:43 IST 2024
```
