// This really helped!
// Source: https://github.com/btracey/mpi/blob/master/examples/bounce/bounce.go

/*
HW04 tests the speed and latency of the underlying network. It sends a suite
of messages between two nodes in the network. This problem must
be run on two nodes.

To run on a single machine:

Then, in two different terminals, run one each of
    bounce -mpi-addr=":5000" -mpi-alladdr=":5000,:5001"
    bounce -mpi-addr=":5001" -mpi-alladdr=":5000,:5001"
*/
package main

import (
	"bytes"
	"encoding/binary"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/btracey/mpi"
)

// length of message. Must be in increasing order
var msgLengths = []int{1e0, 1e1, 1e2, 1e3, 1e4, 1e5, 1e6, 1e7, 1e8, 1e9}

var nRepeats int64 = 100

func main() {
	rand.Seed(time.Now().UnixNano())
	flag.Parse()

	// Initialize MPI
	err := mpi.Init()
	if err != nil {
		log.Fatal("There was an error initializing MPI: ", err)
	}
	defer mpi.Finalize()

	// Second check that things started okay
	rank := mpi.Rank()
	if rank < 0 {
		log.Fatal("Incorrect initialization.")
	}

	size := mpi.Size()
	if size != 2 {
		log.Fatal("ERROR:\tMust use exactly two nodes for this example.\n")
	}
	if rank == 0 {
		fmt.Println("Number of nodes = ", size)
	}

	// Make a random vector of bytes
	maxsize := msgLengths[len(msgLengths)-1]
	message := make([]byte, maxsize)
	for i := 0; i < maxsize/8; i++ {
		v := rand.Int63()
		binary.LittleEndian.PutUint64(message[i*8:], uint64(v))
	}
	receive := make([]byte, maxsize)

	// Do a call and response between the two nodes
	// and record the time. The first nodes will send a message to the next
	// node. The second node will send that message back, which will be received
	// by the first node.
	times := make([]int64, len(msgLengths))
	for i, l := range msgLengths {
		for j := int64(0); j < nRepeats; j++ {
			// Send the byte message
			msg := message[:l]
			rcv := receive[:l]
			start := time.Now()
			if rank == 0 {
				mpi.Send(msg, rank+1, 0) // Send message to next process
			} else if rank == 1 {
				mpi.Receive(&rcv, rank-1, 0) // Get message from process
			}
			if rank == 0 {
				mpi.Receive(&rcv, rank+1, 0)
			} else if rank == 1 {
				mpi.Send(rcv, rank-1, 0)
			}
			times[i] += time.Since(start).Nanoseconds()

			// verify that the received message is the same as the sent message
			if rank == 0 {
				if !bytes.Equal(msg, rcv) {
					log.Fatal("Message not the same. Data loss.")
				}
			}

			// Zero out the buffer so we don't get false positives next time
			for i := range rcv {
				rcv[i] = 0
			}
		}
	}

	// Convert the times to microseconds
	for i := range times {
		times[i] /= time.Microsecond.Nanoseconds()
		times[i] /= nRepeats
	}

	// Have the nodes print their trip time in microseconds
	if rank == 0 {
		for i, message := range msgLengths {
			fmt.Printf("Average %d byte trip time in µs between node %d and %d is : %v\n", message, rank, rank+1, times[i])
		}
	}
}
