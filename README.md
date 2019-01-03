Ctrl-break test
===============

I can't figure out why in a docker container my program can't send windows CTRL-BREAK events to sub-processes.

Setup with

```bash
git clone https://github.com/lox/ctrl-break-test.git
cd ctrl-break-test
go build .
docker build --tag ctrl-break-test .
```

Then run on the host:

```bash
ctrl-break-test.exe
```

Then run in docker:

```bash
docker run --rm ctrl-break-test
```
