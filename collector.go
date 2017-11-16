package main

var WorkQueue = make(chan WorkRequest, 100)

func Collector(r WorkRequest) {
	WorkQueue <- r
	return
}
