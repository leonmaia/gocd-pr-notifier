package main

import (
	"fmt"
)

var WorkQueue = make(chan WorkRequest, 100)

func Collector(r WorkRequest) {

	WorkQueue <- r
	fmt.Println("Work request queued")

	return
}
