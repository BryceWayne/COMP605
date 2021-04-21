package main

import (
	"fmt"
	"strconv"
	"math/rand"
	"os"
	"time"
	"reflect"
)

func main() {
	L := len(os.Args)
	m, n, p, q, err := mapVars(L, os.Args)
	if err != 0 {
	    fmt.Fprintf(os.Stderr, "error: Incorrect command line arguments.\n")
	    os.Exit(1)
	}

	fmt.Println("The product array has dimensions...")
	fmt.Printf("\tC is %dx%d\n", m, q)

	fmt.Println("\nPopulating matrix A.")
	A, _ := createMat(m, n)
	if (m < 5 && q <5) {
		fmt.Println("Matrix A.")
		printMat(m, A)
	}

	fmt.Println("\nPopulating matrix B.")
	B, _ := createMat(p, q)
	if (m <= 5 && q <= 5) {
		fmt.Println("Matrix B.")
		printMat(p, B)
	}

	fmt.Println("\nPerforming row-wise matrix-matrix multiplication AB.")
	C, _ := initMat(m, q)
	startRow := time.Now()
	rowMultMat(m, n, q, A, B, C)
	dtRow := time.Since(startRow)
	fmt.Printf("Time elapsed: %v\n", dtRow)
	if (m <= 5 && q <= 5) {
		fmt.Println("Matrix C.")
		printMat(q, C)
	}

	fmt.Println("\nPerforming column-wise matrix-matrix multiplication AB.")
	D, _ := initMat(m, q)
	startColumn := time.Now()
	colMultMat(m, n, q, A, B, D)
	dtColumn := time.Since(startColumn)
	fmt.Printf("Time elapsed: %v\n", dtColumn)
	if (m <5 && q <5) {
		fmt.Println("Matrix C.")
		printMat(q, D)
	}
	
	if reflect.DeepEqual(C, D) {
		fmt.Println("\nThe two matrices are equivalent.")
	} else {
		fmt.Println("\nThe two matrices are not equivalent.")
	}
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
	} else if (l == 5 || n != p) {
		m, _ := strconv.Atoi(args[1])
		n, _ := strconv.Atoi(args[2])
		p, _ := strconv.Atoi(args[3])
		q, _ := strconv.Atoi(args[4])
		fmt.Println("Creating two arrays, A, B, with dimensions.")
		fmt.Printf("\tA is %dx%d\n\tB is %dx%d\n", m, n, p, q)
		return m, n, p, q, 0
	} else {
		fmt.Println("Incorrect command line arguments.\n")
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
			rows[i*n + j] = float64(rand.Int63()%10)
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
	var total float64 = 0.0
	for i := 0; i < m; i++ {
		for j := 0; j < q; j++ {
			for k := 0; k < n; k++ {
				total += A[i][k] * (B[k][j])
			}
			C[i][j] = total
			total = 0
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