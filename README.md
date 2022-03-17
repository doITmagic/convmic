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
