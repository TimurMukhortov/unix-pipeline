package main

import (
	"sort"
	"strconv"
	"strings"
	"sync"
)

func SingleHash(in chan interface{}, out chan interface{}) {
	for input := range in {
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
			results[1] = DataSignerCrc32(DataSignerMd5(dataStr))
			wg.Done()
		}()
		wg.Wait()
		out <- results[0] + "~" + results[1]
	}
}

func MultiHash(in chan interface{}, out chan interface{}) {
	for input := range in {
		dataStr, ok := input.(string)
		if !ok {
			panic("Could not convert data to string")
		}
		var resultString string
		for number := 0; number < 6; number++ {
			resultString += DataSignerCrc32(strconv.Itoa(number) + dataStr)
		}
		out <- resultString
	}
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
