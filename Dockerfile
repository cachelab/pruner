FROM golang:1.13.4-alpine

MAINTAINER Cache Lab <hello@cachelab.co>

COPY pruner /bin/pruner

USER nobody

ENTRYPOINT ["/bin/pruner"]
