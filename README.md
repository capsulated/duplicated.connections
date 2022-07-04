# Duplicated Connections Checker

Service that compares duplicates in users connections

## Preliminaries

In order to start the service, the following dependencies are required:

* **Go:** service is written in Go. You need to install it on your system
* **dep:** this project uses dep to manage dependencies as well as to provide reproducible builds
* **docker:** this project uses doсker

```
go get -u github.com/golang/dep/cmd/dep
```

### Installing

```
go get -d gilab.com/logiq.tech/foxed.nesthorn
cd $GOPATH/src/gilab.com/logiq.tech/foxed.nesthorn
dep ensure
```

## Run app

```
docker run -p9088:9088 -p6534:6534 -it reindexer/reindexer
go run cmd/service.go
```
## Running the tests

```
docker run -p9088:9088 -p6534:6534 -it reindexer/reindexer
cd tests
ginkgo
```
