package main

import (
    "encoding/gob"
    "fmt"
    "log"
    "net"
    "os"
    "runtime"
    "sync"
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
    ipv4 := os.Args[1]
    fmt.Println("Waiting for connection from:", ipv4)
    // fmt.Println("Listening on port 8080...\nNumber of goroutines: ", runtime.GOMAXPROCS(0))
    ln, err := net.Listen("tcp", ":8080")
    if err != nil {
        log.Fatal("Could not listen:", err)
    }
    conn, err := ln.Accept() // this blocks until connection or error
    if err != nil {
        log.Fatal("Could not accept:", err)
    }
    dec := gob.NewDecoder(conn)
    var data Data
    dec.Decode(&data)
    conn.Close()

    fmt.Println("Received messaage")
    // fmt.Printf("Received : %+v \n", data)
    A := data.A
    B := data.B
    C := data.C
    start := data.RowStart
    end := data.RowEnd
    m := data.M
    n := data.N
    q := data.Q
    goRowMultMat(m, n, q, start, end, A, B, C)

    // fmt.Println("Matrix C.")
    // printMat(m, q, C)

    fmt.Println("Sending result....")
    // Establish connection
    conn, err = net.Dial("tcp", ipv4+":8081")
    if err != nil {
        log.Fatal("Dial:", err)
    }
    message := &Result{}
    message.C = C
    message.Goroutines = runtime.GOMAXPROCS(0)
    encoder := gob.NewEncoder(conn)
    encoder.Encode(message)
    conn.Close()

}

func goRowMultMat(m, n, q, start, end int, A []float64, B []float64, C []float64) {
    var wg sync.WaitGroup
    for i := start; i < end; i++ {
        wg.Add(1)
        go func(i int) {
            for j := 0; j < q; j++ {
                C[i*q+j] = 0
                for k := 0; k < n; k++ {
                    C[i*q+j] = C[i*q+j] + A[i*n+k]*B[k*q+j]
                }
            }
            wg.Done()
        }(i)
    }
    wg.Wait()
}

func printMat(m, n int, M []float64) {
    for i := 0; i < m; i++ {
        for j := 0; j < n; j++ {
            fmt.Printf("%v ", M[i*n+j])
        }
        fmt.Print("\n")
    }
}
