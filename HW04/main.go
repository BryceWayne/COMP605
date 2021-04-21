package main

import (
	"fmt"
	"math"
	"os"
	"strconv"
	"sync"
	"time"
)

func f(x float64) float64 {
	return math.Sin(x) * math.Cos(2*x)
}

func main() {
	if len(os.Args) < 4 {
		fmt.Println("\n\t\t ~~~~ HALP ~~~~ \n")
		fmt.Println("This program is intended to perform integration via the trapezoid method.")
		fmt.Println("\t Current function is f(x)=sin(x)*cos(2x)\n")
		fmt.Println("TEMPLATE: go run main.go a b n\n")
		fmt.Println("EXAMPLE:  go run main.go 0 1 1000\n------->  Total Time: 4.9891ms, Integral: 0.354037\n")
		fmt.Println("PARAMETERS:")
		fmt.Println("\t a: (FLOAT) Starting value")
		fmt.Println("\t b: (FLOAT) Ending value")
		fmt.Println("\t n:  (INT)  Number of trapezoids to create")
		return
	}
	var wg sync.WaitGroup

	a, _ := strconv.ParseFloat(os.Args[1], 64)
	b, _ := strconv.ParseFloat(os.Args[2], 64)
	n, _ := strconv.Atoi(os.Args[3])
	var in = make([]float64, n, n)

	dx := (b - a) / float64(n)
	approx := 0.5 * (f(a) + f(b))
	// fmt.Printf("h: %f, approx: %f\n", h, approx)
	start := time.Now()
	for i := 0; i < n; i++ {
		wg.Add(1)
		x_i := a + float64(i)*dx
		go func(i int, x float64, in []float64) {
			y := f(x)
			in[i] = y
			wg.Done()
		}(i, x_i, in)
	}
	wg.Wait()
	dt := time.Since(start)
	for _, v := range in {
		approx += v
	}
	approx = dx * approx
	fmt.Printf("Total Time: %v, Integral: %f\n", dt, approx)
}
