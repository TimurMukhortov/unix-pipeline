package main

import (
	"runtime"
)

func SingleHash(in chan interface{}, out chan interface{}) {

}

func MultiHash(in chan interface{}, out chan interface{}) {

}

func CombineResults(in chan interface{}, out chan interface{}) {

}

func ExecutePipeline(jobs ...job) {
	inCh := make(chan interface{})
	outCh := make(chan interface{})
	for _, job := range jobs {
		go job(inCh, outCh)
		inCh = outCh
		outCh = make(chan interface{})
		runtime.Gosched()
	}
}

func main() {
}
