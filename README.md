# GRPC Micrososervice for currency conversion
 just a demo application, not for production


## Test files 
are in bitbucket.org/doitmagic/convmic/src/server/testing
```
cd src/server/testing
go test
```

## Profiling
 profiling file is generated only if you start server application with flag, in server dir  
  ```
  -cpuprofile true 
  ```

## Protocol Buffer Compiler  command
```protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative proto/currencies.proto```

## Concurrents writes
- The currencies variable from appcontext is read and write by concurrent goroutines.
- We chose to use sync.Map {} because: ```The Map type is optimized for two common use cases: (1) when the entry for a given key is only ever written once but read many times, as in caches that only grow, or (2) when multiple goroutines read, write, and overwrite entries for disjoint sets of keys. In these two cases, use of a Map may significantly reduce lock contention compared to a Go map paired with a separate Mutex or RWMutex.```

## How to test the program:

### The GRPC Server :
```
  cd src/server/cmd

  # start the server
  go run main

  # to start with profiling 
  go run main.go -cpuprofile yes 
``` 

### The GRPC client :
```
  cd src/client/cmd
  
  # list the currencies with pagination `-c` is command, `-p` page, `-pp` records per page  
  go run main.go -c list -p 1 -pp 10

  # convert one or multiple currencies  `-c` is command, `-p` page, `-pp` records per page 
  # the currecies to convert is hardcoded in main.go 
  go run main.go -c convert 

``` 