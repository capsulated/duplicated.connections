# Foxed Microservice [demo for portfolio]

Service that compares duplicates in users connections

## Preliminaries

In order to start the service, the following dependencies are required:

* **Go:** service is written in Go. You need to install it on your system
* **dep:** this project uses dep to manage dependencies as well as to provide reproducible builds

```
go get -u github.com/golang/dep/cmd/dep
```

### Installing

```
go get -d gilab.com/logiq.tech/foxed.nesthorn
cd $GOPATH/src/gilab.com/logiq.tech/foxed.nesthorn
dep ensure

```
