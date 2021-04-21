package main

import (
	"fmt"
	"math"
	"math/rand"
	"os"
	"strconv"
	"sync"
	"time"
)

func main() {
	N, _ := strconv.Atoi(os.Args[1])
	var in = make([]int, N, N)
	var wg sync.WaitGroup
	start := time.Now()
	for i := 0; i < N; i++ {
		wg.Add(1)
		go func(i int, in []int) {
			x := 2*rand.Float64() - 1.0
			y := 2*rand.Float64() - 1.0
			value := math.Sqrt(x*x + y*y)
			if value <= 1.0 {
				in[i] = 1
			}
			wg.Done()
		}(i, in)
	}
	wg.Wait()
	dt := time.Since(start)
	var summand int = 0
	for _, v := range in {
		summand += v
	}
	fmt.Printf("Total Time: %v, In: %d, Total: %d, Approx. Pi: %f\n", dt, summand, N, 4*float64(summand)/float64(N))
}
