# Foxed Microservice

Service that compares duplicates in users connections

## Preliminaries

In order to start the service, the following dependencies are required:

* **Go:** service is written in Go. You need to install it on your system
* **dep:** this project uses dep to manage dependencies as well as to provide reproducible builds
* **reindexer:** this project uses reindexer inmemory db

```
go get -u github.com/golang/dep/cmd/dep
go get -a github.com/restream/reindexer
bash $GOPATH/src/github.com/restream/reindexer/dependencies.sh
go generate github.com/restream/reindexer/bindings/builtin
```

### Installing

```
go get -d gilab.com/logiq.tech/foxed.nesthorn
cd $GOPATH/src/gilab.com/logiq.tech/foxed.nesthorn
dep ensure

```

## Running the tests

```
cd tests
ginkgo
```

## PS

В докер контейнер не стал класть, времени не хватило.