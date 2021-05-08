package main

import (
    "encoding/gob"
    "fmt"
    "log"
    "math/rand"
    "net"
    "os"
    "runtime"
    "strconv"
    "sync"
    "time"
)

type Data struct {
    M, N, Q          int
    RowStart, RowEnd int
    A, B, C          []float64
}

type Result struct {
    Goroutines int
    C          []float64
}

func main() {
    // Get vars
    m, n, p, q, ipv4, errInt := mapVars(os.Args)
    if errInt != 0 {
        return
    }
    fmt.Println("Connecting to: ", ipv4)
    fmt.Println("Starting client....\nNumber of goroutines: ", runtime.GOMAXPROCS(0))
    // Establish connection
    conn, err := net.Dial("tcp", ipv4+":8080")
    if err != nil {
        log.Fatal("Dial:", err)
    }

    // Init params
    fmt.Println("\nPopulating matrix A.")
    A := createMat(m, n)

    fmt.Println("Populating matrix B.")
    B := createMat(p, q)

    C := initMat(m, q)

    // Compute message info
    message := &Data{}
    message.M = m
    message.N = n
    message.Q = q
    rowStart := m / 2
    message.RowStart = rowStart
    message.RowEnd = m
    message.A = A
    message.B = B
    message.C = C
    // Encode data struct and send instructions
    encoder := gob.NewEncoder(conn)
    encoder.Encode(message)
    conn.Close()
    fmt.Println("\nMessage sent. Starting matrix-matrix multiplication.\n")

    startRow := time.Now()
    // Do some stuff in a goroutine
    var wg sync.WaitGroup
    wg.Add(1)
    go goRowMultMat(m, n, q, 0, rowStart, A, B, C, &wg)

    // Wait for other process to signal done
    fmt.Println("Waiting for results...")
    ln, err := net.Listen("tcp", ":8081")
    if err != nil {
        log.Fatal("Could not listen:", err)
    }
    conn, err = ln.Accept() // this blocks until connection or error
    if err != nil {
        log.Fatal("Could not accept:", err)
    }
    dec := gob.NewDecoder(conn)
    result := &Result{}
    dec.Decode(result)
    conn.Close()
    wg.Wait()

    _ = combineMatrix(m, q, 0, rowStart, C, result.C)
    dtRow := time.Since(startRow)
    if m <= 5 && q <= 5 {
        fmt.Println("Matrix C.")
        printMat(m, q, C)
    }

    fmt.Printf("Time elapsed: %v \n", dtRow)
    fmt.Println("Total number of goroutines used:", runtime.GOMAXPROCS(0)+result.Goroutines)
}

func mapVars(args []string) (m int, n int, p int, q int, ipv4 string, err int) {
    l := len(args)
    if l == 3 {
        m, _ := strconv.Atoi(args[1])
        n, _ := strconv.Atoi(args[1])
        p, _ := strconv.Atoi(args[1])
        q, _ := strconv.Atoi(args[1])
        ipv4 := args[2]
        fmt.Printf("Creating two arrays, A, B, with square dimensions.\n")
        fmt.Printf("\tA is %dx%d\n\tB is %dx%d\n", m, n, p, q)
        return m, n, p, q, ipv4, 0
    } else if l == 6 || args[2] != args[3] {
        m, _ := strconv.Atoi(args[1])
        n, _ := strconv.Atoi(args[2])
        p, _ := strconv.Atoi(args[3])
        q, _ := strconv.Atoi(args[4])
        ipv4 := args[5]
        fmt.Println("Creating two arrays, A, B, with dimensions.")
        fmt.Printf("\tA is %dx%d\n\tB is %dx%d\n", m, n, p, q)
        return m, n, p, q, ipv4, 0
    } else {
        fmt.Println("````````````````````````````````````````````````````````````````````````````````\n\n")
        fmt.Println("\tALERT:Incorrect number of input arguments.\n\t                        Exiting.\n\n")
        fmt.Println("````````````````````````````````````````````````````````````````````````````````\n\n")
        fmt.Println("\tUsage:\n")
        fmt.Println("\t$ args rowsA columnsA rowsB columnsB IPv4\n")
        fmt.Println("\trowsA: The number of rows in Matrix A.\n")
        fmt.Println("\tcolumnsA: The number of columns in Matrix A.\n")
        fmt.Println("\trowsB: The number of rows in Matrix B.\n")
        fmt.Println("\tcolumnsB: The number of columns in Matrix B.\n")
        fmt.Println("\tIPv4: The IPv4 address of the current machine.\n")
        return 0, 0, 0, 0, "", 1
    }
}

func initMat(m int, n int) []float64 {
    M := make([]float64, m*n)
    return M
}

func createMat(m int, n int) []float64 {
    M := make([]float64, m*n)

    for i := 0; i < m; i++ {
        for j := 0; j < n; j++ {
            M[i*n+j] = rand.Float64()
        }
    }
    return M
}

func printMat(m, n int, M []float64) {
    for i := 0; i < m; i++ {
        for j := 0; j < n; j++ {
            fmt.Printf("%v ", M[i*n+j])
        }
        fmt.Print("\n")
    }
}

func goRowMultMat(m, n, q, start, end int, A []float64, B []float64, C []float64, wg *sync.WaitGroup) {
    defer wg.Done()
    var wg_internal sync.WaitGroup
    for i := start; i < end; i++ {
        wg_internal.Add(1)
        go func(i int) {
            for j := 0; j < q; j++ {
                C[i*q+j] = 0
                for k := 0; k < n; k++ {
                    C[i*q+j] = C[i*q+j] + A[i*n+k]*B[k*q+j]
                }
            }
            wg_internal.Done()
        }(i)
    }
    wg_internal.Wait()
}

func combineMatrix(m, q, start, stop int, A, B []float64) []float64 {
    // performs matrix addition: A + B => move B elements into A
    for i := start; i < stop; i++ {
        for j := 0; j < q; j++ {
            B[i*q+j] = A[i*q+j]
        }
    }
    A = B
    return A
}
