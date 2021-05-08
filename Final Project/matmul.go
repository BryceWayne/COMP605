package main

import (
	"fmt"
	"math/rand"
	"os"
	"runtime"
	"strconv"
	"sync"
	"time"
)

func main() {
	L := len(os.Args)
	m, n, p, q, err := mapVars(L, os.Args)
	if err != 0 {
		return
	}

	fmt.Println("The product array has dimensions...")
	fmt.Printf("\tC is %dx%d\n", m, q)

	fmt.Println("\nPopulating matrix A.")
	A, _ := createMat(m, n)
	if m < 5 && n < 5 {
		fmt.Println("Matrix A.")
		printMat(m, A)
	}

	fmt.Println("Populating matrix B.")
	B, _ := createMat(p, q)
	if p <= 5 && q <= 5 {
		fmt.Println("Matrix B.")
		printMat(p, B)
	}

	fmt.Println("\nPerforming row-wise matrix-matrix multiplication AB.")
	C, _ := initMat(m, q)
	startRow := time.Now()
	rowMultMat(m, n, q, A, B, C)
	dtRow := time.Since(startRow)
	fmt.Printf("Time elapsed: %v\n", dtRow)

	fmt.Printf("\nPerforming row-wise matrix-matrix multiplication AB using %d goroutines.\n", runtime.GOMAXPROCS(0))
	E, _ := initMat(m, q)

	startGo := time.Now()
	goRowMultMat(m, n, q, A, B, E)
	dtGo := time.Since(startGo)

	fmt.Printf("Time elapsed: %v\n", dtGo)
}

func mapVars(l int, args []string) (m int, n int, p int, q int, err int) {
	if l == 2 {
		m, _ := strconv.Atoi(args[1])
		n, _ := strconv.Atoi(args[1])
		p, _ := strconv.Atoi(args[1])
		q, _ := strconv.Atoi(args[1])
		fmt.Printf("Creating two arrays, A, B, with square dimensions.\n")
		fmt.Printf("\tA is %dx%d\n\tB is %dx%d\n", m, n, p, q)
		return m, n, p, q, 0
	} else if l == 5 || n != p {
		m, _ := strconv.Atoi(args[1])
		n, _ := strconv.Atoi(args[2])
		p, _ := strconv.Atoi(args[3])
		q, _ := strconv.Atoi(args[4])
		fmt.Println("Creating two arrays, A, B, with dimensions.")
		fmt.Printf("\tA is %dx%d\n\tB is %dx%d\n", m, n, p, q)
		return m, n, p, q, 0
	} else {
		fmt.Println("````````````````````````````````````````````````````````````````````````````````\n\n")
		fmt.Println("\tALERT:Incorrect number of input arguments.\n\t      Exiting.\n\n")
		fmt.Println("````````````````````````````````````````````````````````````````````````````````\n\n")
		fmt.Println("\tUsage:\n")
		fmt.Println("\t$ args rowsA columnsA rowsB columnsB\n")
		fmt.Println("\trowsA: The number of rows in Matrix A.\n")
		fmt.Println("\tcolumnsA: The number of columns in Matrix A.\n")
		fmt.Println("\trowsB: The number of rows in Matrix B.\n")
		fmt.Println("\tcolumnsB: The number of columns in Matrix B.\n")
		return 0, 0, 0, 0, 1
	}
}

func initMat(m int, n int) (M [][]float64, rows []float64) {
	M = make([][]float64, m)
	rows = make([]float64, n*m)
	for i := 0; i < m; i++ {
		M[i] = rows[i*n : (i+1)*n]
	}
	return M, rows
}

func createMat(m int, n int) (M [][]float64, rows []float64) {
	M = make([][]float64, m)
	rows = make([]float64, n*m)
	for i := 0; i < m; i++ {
		for j := 0; j < n; j++ {
			rows[i*n+j] = rand.Float64()
		}
		M[i] = rows[i*n : (i+1)*n]
	}
	return M, rows
}

func printMat(row int, M [][]float64) {
	for i := 0; i < row; i++ {
		fmt.Printf("%v\n", M[i])
	}
}

func rowMultMat(m int, n int, q int, A [][]float64, B [][]float64, C [][]float64) {
	for i := 0; i < m; i++ {
		for j := 0; j < q; j++ {
			C[i][j] = 0
			for k := 0; k < n; k++ {
				C[i][j] = C[i][j] + A[i][k]*(B[k][j])
			}
		}
	}
}

func colMultMat(m int, n int, q int, A [][]float64, B [][]float64, C [][]float64) {
	for j := 0; j < q; j++ {
		for i := 0; i < m; i++ {
			C[i][j] = 0
		}
		for k := 0; k < n; k++ {
			for i := 0; i < m; i++ {
				C[i][j] += A[i][k] * (B[k][j])
			}
		}
	}
}

func goRowMultMat(m int, n int, q int, A [][]float64, B [][]float64, C [][]float64) {
	var wg sync.WaitGroup
	for i := 0; i < m; i++ {
		wg.Add(1)
		go func(i int) {
			for j := 0; j < q; j++ {
				C[i][j] = 0
				for k := 0; k < n; k++ {
					C[i][j] = C[i][j] + A[i][k]*(B[k][j])
				}
			}
			wg.Done()
		}(i)
	}
	wg.Wait()
}
