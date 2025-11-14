## Remote Execution

The Remote Execution package provides secure remote command execution capabilities, allowing clients to execute commands on remote servers over a network connection.

## Table of Contents

1. [What is Remote Execution?](#what-is-remote-execution)
2. [How Does It Work?](#how-does-it-work)
3. [Understanding the Code](#understanding-the-code)
4. [Usage Examples](#usage-examples)
5. [Security Considerations](#security-considerations)
6. [Further Reading](#further-reading)

### What is Remote Execution?

Remote Execution allows running commands and scripts on remote systems over a network. It's a fundamental capability for system administration, automation, and distributed computing.

**Key Features**:
- Execute commands on remote systems
- Authentication token support
- Command whitelisting for security
- Synchronous and asynchronous execution
- Batch command execution
- Interactive command shell
- Execution tracking and management

**Common Use Cases**:
- System administration and maintenance
- Deployment automation
- Distributed task execution
- Remote monitoring and diagnostics
- CI/CD pipelines
- Configuration management

### How Does It Work?

The remote execution system uses a client-server architecture:

1. **Server Setup**: Server starts and listens for connections
2. **Client Connection**: Client connects to server via TCP
3. **Authentication**: Client provides authentication token (if required)
4. **Command Execution**: Client sends command request
5. **Server Processing**: Server validates, executes command
6. **Response**: Server returns output, error, and exit code

**Execution Flow**:
```
Client                          Server
  |                               |
  |--- Connect ------------------->|
  |                               |
  |--- CommandRequest ------------>|
  |   {cmd, args, token}          |
  |                               |--- Authenticate
  |                               |--- Validate Command
  |                               |--- Execute Command
  |                               |
  |<-- CommandResponse ------------|
  |   {output, error, exitCode}   |
```

### Understanding the Code

#### Server Components:

- `Server`: Manages remote execution server
- `ExecutionContext`: Tracks active command executions
- `CommandRequest`: Client request structure
- `CommandResponse`: Server response with results
- `handleConnection`: Processes client connections
- `executeCommand`: Runs commands and captures output

**Security Features**:
- **Authentication**: Token-based authentication
- **Command Whitelisting**: Only allowed commands can execute
- **Execution Tracking**: Monitor active executions
- **Process Management**: Kill running commands if needed

#### Client Components:

- `Client`: Manages connection to remote execution server
- `Execute()`: Synchronous command execution
- `ExecuteAsync()`: Asynchronous command execution
- `ExecuteWithTimeout()`: Command execution with timeout
- `ExecuteScript()`: Execute shell scripts
- `ExecuteBatch()`: Execute multiple commands in sequence
- `RunInteractive()`: Interactive command shell

### Usage Examples

#### Starting a Remote Execution Server:
```go
// Define allowed commands for security
allowedCmds := []string{"ls", "pwd", "whoami", "date", "echo"}

// Create server with authentication
authToken := "your-secret-token-here"
server := remoteexec.NewServer(":9000", authToken, allowedCmds)

if err := server.Start(); err != nil {
    log.Fatal(err)
}
```

#### Connecting a Client and Executing Commands:
```go
client := remoteexec.NewClient("server.example.com:9000", "your-secret-token-here")

// Connect to server
if err := client.Connect(); err != nil {
    log.Fatal(err)
}
defer client.Disconnect()

// Execute a command
response, err := client.Execute("ls", "-la", "/home")
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Exit Code: %d\n", response.ExitCode)
fmt.Printf("Output:\n%s\n", response.Output)
```

#### Asynchronous Execution:
```go
// Execute command asynchronously
responseChan, err := client.ExecuteAsync("sleep", "5")
if err != nil {
    log.Fatal(err)
}

// Do other work while command executes
fmt.Println("Command running in background...")

// Wait for result
response := <-responseChan
fmt.Printf("Command completed: %v\n", response.Success)
```

#### Execution with Timeout:
```go
// Execute with 10-second timeout
response, err := client.ExecuteWithTimeout(10*time.Second, "long-running-command")
if err != nil {
    log.Printf("Command timed out: %v\n", err)
} else {
    fmt.Printf("Output: %s\n", response.Output)
}
```

#### Executing Shell Scripts:
```go
script := `
#!/bin/bash
echo "Starting backup..."
tar -czf backup.tar.gz /data
echo "Backup complete"
`

response, err := client.ExecuteScript(script)
if err != nil {
    log.Fatal(err)
}

fmt.Println(response.Output)
```

#### Batch Execution:
```go
commands := [][]string{
    {"mkdir", "-p", "/tmp/mydir"},
    {"cd", "/tmp/mydir"},
    {"touch", "file1.txt"},
    {"ls", "-la"},
}

responses, err := client.ExecuteBatch(commands)
if err != nil {
    log.Printf("Batch failed: %v\n", err)
}

for i, resp := range responses {
    fmt.Printf("Command %d: %s\n", i, resp.Output)
}
```

#### Interactive Mode:
```go
client := remoteexec.NewClient("localhost:9000", "token")
client.Connect()
defer client.Disconnect()

// Start interactive shell
client.RunInteractive()
```

#### Server Management:
```go
// Check active executions
count := server.GetActiveExecutions()
fmt.Printf("Active executions: %d\n", count)

// Kill a specific execution
err := server.KillExecution("execution-id")

// Dynamically add allowed command
server.AddAllowedCommand("git")

// Check if command is allowed
if server.IsCommandAllowed("rm") {
    fmt.Println("rm is allowed")
}
```

### Security Considerations

**⚠️ IMPORTANT SECURITY WARNINGS**:

1. **Authentication**: Always use authentication tokens in production
2. **Command Whitelisting**: Never allow all commands (`allowedCmds` should be restrictive)
3. **Network Security**: Use in trusted networks or over encrypted tunnels (VPN/TLS)
4. **Input Validation**: Commands are executed with shell access
5. **Access Control**: Limit which users/systems can connect
6. **Logging**: Log all executed commands for audit trails

**Best Practices**:
- Run server with minimal privileges (non-root user)
- Use strong, randomly generated tokens
- Implement IP whitelisting
- Add TLS encryption for network transport
- Monitor and rate-limit command executions
- Implement command output size limits
- Set execution timeouts server-side

**Production Considerations**:
```go
// Use environment variables for tokens
authToken := os.Getenv("REMOTE_EXEC_TOKEN")

// Strict command whitelist
allowedCmds := []string{
    "systemctl status",
    "journalctl",
    "df",
    "free",
    "uptime",
}

// No wildcard command allowing
server := remoteexec.NewServer(":9000", authToken, allowedCmds)
```

### Protocol Details

**Request Format**:
```json
{
  "id": "uuid-v4",
  "command": "ls",
  "args": ["-la", "/home"],
  "token": "auth-token-here"
}
```

**Response Format**:
```json
{
  "id": "uuid-v4",
  "success": true,
  "output": "total 8\ndrwxr-xr-x 2 user user 4096...",
  "error": "",
  "exit_code": 0,
  "duration": "15.234ms"
}
```

### Further Reading

- [SSH Remote Execution](https://man.openbsd.org/ssh)
- [PowerShell Remoting](https://docs.microsoft.com/en-us/powershell/scripting/learn/remoting/running-remote-commands)
- [Ansible Architecture](https://docs.ansible.com/ansible/latest/dev_guide/overview_architecture.html)
- [Remote Procedure Calls](https://en.wikipedia.org/wiki/Remote_procedure_call)
- [OWASP Secure Coding Practices](https://owasp.org/www-project-secure-coding-practices-quick-reference-guide/)
