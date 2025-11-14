## nslookup Implementation

This package implements DNS lookup functionality similar to the `nslookup` command-line tool.

### Features

- A record lookup (hostname to IP)
- Reverse DNS lookup (IP to hostname)
- MX record lookup (mail servers)
- NS record lookup (name servers)
- TXT record lookup
- CNAME record lookup
- Lookup with timeout support
- Comprehensive DNS queries

### Usage

#### Basic Hostname Lookup

```go
package main

import (
	"fmt"
	"log"
	"networkprogramming/nslookup"
)

func main() {
	result, err := nslookup.Lookup("google.com")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Print(nslookup.FormatResult(result))
}
```

#### Reverse DNS Lookup

```go
package main

import (
	"fmt"
	"log"
	"networkprogramming/nslookup"
)

func main() {
	result, err := nslookup.LookupIP("8.8.8.8")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Print(nslookup.FormatResult(result))
}
```

#### MX Record Lookup

```go
package main

import (
	"fmt"
	"log"
	"networkprogramming/nslookup"
)

func main() {
	result, err := nslookup.LookupMX("gmail.com")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Print(nslookup.FormatResult(result))
}
```

#### Lookup All Records

```go
package main

import (
	"fmt"
	"networkprogramming/nslookup"
)

func main() {
	results := nslookup.LookupAll("example.com")

	for recordType, result := range results {
		fmt.Printf("=== %s Records ===\n", recordType)
		fmt.Print(nslookup.FormatResult(result))
		fmt.Println()
	}
}
```

#### Lookup with Timeout

```go
package main

import (
	"fmt"
	"log"
	"time"
	"networkprogramming/nslookup"
)

func main() {
	result, err := nslookup.LookupWithTimeout("google.com", 5*time.Second)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Print(nslookup.FormatResult(result))
}
```

### Record Types

- **A**: IPv4 address records
- **AAAA**: IPv6 address records
- **CNAME**: Canonical name (alias) records
- **MX**: Mail exchange records
- **NS**: Name server records
- **TXT**: Text records (SPF, DKIM, etc.)

### Output Format

```
Query: google.com
Addresses:
  142.250.185.46
  2607:f8b0:4004:c07::71

Query: gmail.com
MX Records:
  gmail-smtp-in.l.google.com (priority: 5)
  alt1.gmail-smtp-in.l.google.com (priority: 10)
```

### How it Works

1. Uses Go's `net` package for DNS resolution
2. Queries system's configured DNS servers
3. Supports multiple record types
4. Returns structured results with all found records
5. Handles timeouts and errors gracefully

### Use Cases

- DNS troubleshooting
- Verify domain configuration
- Check mail server setup (MX records)
- Find authoritative name servers
- Verify DNS propagation
- Security research (TXT records, SPF)
- Network diagnostics

### Common Operations

#### Check if domain exists
```go
result, err := nslookup.Lookup("example.com")
if err != nil {
    fmt.Println("Domain not found or DNS error")
}
```

#### Find mail servers
```go
result, _ := nslookup.LookupMX("domain.com")
for _, mx := range result.MXRecords {
    fmt.Printf("Mail server: %s\n", mx.Host)
}
```

#### Check name servers
```go
result, _ := nslookup.LookupNS("domain.com")
for _, ns := range result.NSRecords {
    fmt.Printf("Name server: %s\n", ns)
}
```

### Limitations

- Uses system's default DNS resolver
- Does not support custom DNS server specification in basic mode
- DNSSEC validation not implemented
- Limited to standard Go DNS capabilities
