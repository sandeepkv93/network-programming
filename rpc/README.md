## Remote Procedure Call (RPC)

Go's RPC package implementation providing remote method invocation over TCP connections.

## Features

- Remote method invocation
- Built-in Go RPC protocol
- Synchronous calls
- Service registration
- Type-safe method calls

## Usage

### Server
```go
server := rpc.NewServer(":1234")
server.Start()
```

### Client
```go
client := rpc.NewClient("localhost:1234")
client.Connect()
result, err := client.Add(5, 3)
fmt.Printf("5 + 3 = %d\n", result)
```

## Built-in MathService

- Add(a, b int) int
- Subtract(a, b int) int
- Multiply(a, b int) int
- Divide(a, b int) int
