go build .
docker build --tag ctrl-break-test .
docker run --rm ctrl-break-test