package main

import (
	"sort"
	"strconv"
	"strings"
	"sync"
)

func SingleHash(in chan interface{}, out chan interface{}) {
	mu := &sync.Mutex{}
	loopWg := sync.WaitGroup{}
	for input := range in {
		input := input
		loopWg.Add(1)
		go func() {
			dataInt, ok := input.(int)
			if !ok {
				panic("Could not convert data to int")
			}
			dataStr := strconv.Itoa(dataInt)
			wg := sync.WaitGroup{}
			wg.Add(2)
			results := make([]string, 2)
			go func() {
				results[0] = DataSignerCrc32(dataStr)
				wg.Done()
			}()
			go func() {
				var dataSignerResult string
				mu.Lock()
				dataSignerResult = DataSignerMd5(dataStr)
				mu.Unlock()
				results[1] = DataSignerCrc32(dataSignerResult)
				wg.Done()
			}()
			wg.Wait()
			out <- results[0] + "~" + results[1]
			loopWg.Done()
		}()
	}
	loopWg.Wait()
}

func MultiHash(in chan interface{}, out chan interface{}) {
	loopWg := sync.WaitGroup{}
	for input := range in {
		input := input
		loopWg.Add(1)
		go func() {
			dataStr, ok := input.(string)
			if !ok {
				panic("Could not convert data to string")
			}
			results := make([]string, 6)
			wg := sync.WaitGroup{}
			for number := 0; number < 6; number++ {
				wg.Add(1)
				number := number
				go func() {
					results[number] = DataSignerCrc32(strconv.Itoa(number) + dataStr)
					wg.Done()
				}()
			}
			wg.Wait()
			out <- strings.Join(results, "")
			loopWg.Done()
		}()
	}
	loopWg.Wait()
}

func CombineResults(in chan interface{}, out chan interface{}) {
	values := []string{}
	for input := range in {
		dataStr, ok := input.(string)
		if !ok {
			panic("Could not convert data to string")
		}
		values = append(values, dataStr)
	}
	sort.Strings(values)
	out <- strings.Join(values, "_")
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
