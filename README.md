# Pruner

Monitors elasticsearch time based indices and will prune on a configurable setup.

[![CircleCI](https://circleci.com/gh/cachelab/pruner.svg?style=svg)](https://circleci.com/gh/cachelab/pruner)

## Usage

This task is configured by the following environment variables:

```bash
ELASTICSEARCH_URL # the url to the elasticsearch cluster
DELAY             # how long to delay until querying the indices again
PREFIXES          # csv separated list of indices to prune
DAYS_TO_KEEP      # how many days of logs you wish to keep
MAX_RETRIES       # how many times the client will try to connect to Elasticsearch
RUN_ONCE          # used for unit testing or to manually prune once
```

## Contributing

* `make run` - runs the pruner in a docker container
* `make build` - builds your pruner docker container
* `make vet` - go fmt and vet code
* `make test` - run unit tests

Before you submit a pull request please update the semantic version inside of
`main.go` with what you feel is appropriate and then edit the `CHANGELOG.md` with
your changes and follow a similar structure to what is there.
