package main

func SingleHash(in chan interface{}, out chan interface{}) {

}

func MultiHash(in chan interface{}, out chan interface{}) {

}

func CombineResults(in chan interface{}, out chan interface{}) {

}

func ExecutePipeline(jobs ...job) {
	var (
		inCh  chan interface{}
		outCh chan interface{}
	)
	outCh = make(chan interface{})
	close(outCh)
	for _, j := range jobs {
		inCh = outCh
		outCh = make(chan interface{})
		go func(in chan interface{}, out chan interface{}, j job) {
			j(in, out)
			close(out)
		}(inCh, outCh, j)
	}
	for {
		select {
		case _, ok := <-outCh:
			if !ok {
				return
			}
		}
	}
}

func main() {
}
