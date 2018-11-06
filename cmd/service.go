package main

import (
	"github.com/logiqone/foxed.nesthorn/workers"
)

func main() {

	fox := new(workers.Foxer)
	if err := fox.Init(); err != nil {
		panic(err)
	}

	if err := fox.InitDataFill(); err != nil {
		panic(err)
	}

	if err := fox.Start(); err != nil {
		println("Error in ListenAndServe: %s", err)
	}
}