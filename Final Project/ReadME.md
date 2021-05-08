# Read THIS
 CAT this!

# HowTo
This readME is designed to help you run this repo.

   1. Download Go
   2. Install Go
   3. Profit

   *Just kidding.*

The three main files for computation are `matmul.go`, `client.go` and `server.go`. The scripts, along with their binaries in the `./bin` folder, will allow you to verify my results. Each machine is different so please understand that the timings are hardware dependent.

 1. To execute the binary for `matmul.go` please see below...
  > `cd ./bin`
  
  > `./matmul n`
  
  > Example: `./matmul 1000`

 2. To execute a distributed matrix-matrix computation please open two terminals.
  > In Terminal one, `rsh node10`. This will connect to node 10 in the network.
  
  > Next, navigate to the bin directory for this project
  
  > `./server 10.0.1.12` will create a server to listen for a message from node12.
  
  > In Terminal two, `rsh node12`. This will connect to node 12 in the network.
  
  > Navigate to the bin directory for the project.
  
  > `./client 1000 10.0.1.10` will launch a client, send computation directions to node10, then perform its half of the work, wait for node10 to send back its computation and finally combine the results.
