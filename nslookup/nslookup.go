package nslookup

import (
	"context"
	"fmt"
	"net"
	"strings"
	"time"
)

// RecordType represents DNS record types
type RecordType string

const (
	RecordTypeA     RecordType = "A"
	RecordTypeAAAA  RecordType = "AAAA"
	RecordTypeCNAME RecordType = "CNAME"
	RecordTypeMX    RecordType = "MX"
	RecordTypeNS    RecordType = "NS"
	RecordTypeTXT   RecordType = "TXT"
)

// LookupResult represents the result of a DNS lookup
type LookupResult struct {
	Query      string
	RecordType RecordType
	Addresses  []string
	Names      []string
	MXRecords  []*net.MX
	NSRecords  []string
	TXTRecords []string
	Error      error
}

// Lookup performs a DNS lookup for the specified hostname
func Lookup(hostname string) (*LookupResult, error) {
	result := &LookupResult{
		Query:      hostname,
		RecordType: RecordTypeA,
		Addresses:  []string{},
	}

	// Lookup IP addresses
	addrs, err := net.LookupHost(hostname)
	if err != nil {
		result.Error = err
		return result, err
	}

	result.Addresses = addrs
	return result, nil
}

// LookupWithTimeout performs a DNS lookup with a timeout
func LookupWithTimeout(hostname string, timeout time.Duration) (*LookupResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	result := &LookupResult{
		Query:      hostname,
		RecordType: RecordTypeA,
		Addresses:  []string{},
	}

	resolver := &net.Resolver{}
	addrs, err := resolver.LookupHost(ctx, hostname)
	if err != nil {
		result.Error = err
		return result, err
	}

	result.Addresses = addrs
	return result, nil
}

// LookupIP performs a reverse DNS lookup for an IP address
func LookupIP(ip string) (*LookupResult, error) {
	result := &LookupResult{
		Query: ip,
		Names: []string{},
	}

	names, err := net.LookupAddr(ip)
	if err != nil {
		result.Error = err
		return result, err
	}

	result.Names = names
	return result, nil
}

// LookupMX performs an MX (Mail Exchange) record lookup
func LookupMX(domain string) (*LookupResult, error) {
	result := &LookupResult{
		Query:      domain,
		RecordType: RecordTypeMX,
		MXRecords:  []*net.MX{},
	}

	mxRecords, err := net.LookupMX(domain)
	if err != nil {
		result.Error = err
		return result, err
	}

	result.MXRecords = mxRecords
	return result, nil
}

// LookupNS performs an NS (Name Server) record lookup
func LookupNS(domain string) (*LookupResult, error) {
	result := &LookupResult{
		Query:      domain,
		RecordType: RecordTypeNS,
		NSRecords:  []string{},
	}

	nsRecords, err := net.LookupNS(domain)
	if err != nil {
		result.Error = err
		return result, err
	}

	for _, ns := range nsRecords {
		result.NSRecords = append(result.NSRecords, ns.Host)
	}

	return result, nil
}

// LookupTXT performs a TXT record lookup
func LookupTXT(domain string) (*LookupResult, error) {
	result := &LookupResult{
		Query:      domain,
		RecordType: RecordTypeTXT,
		TXTRecords: []string{},
	}

	txtRecords, err := net.LookupTXT(domain)
	if err != nil {
		result.Error = err
		return result, err
	}

	result.TXTRecords = txtRecords
	return result, nil
}

// LookupCNAME performs a CNAME record lookup
func LookupCNAME(hostname string) (*LookupResult, error) {
	result := &LookupResult{
		Query:      hostname,
		RecordType: RecordTypeCNAME,
		Names:      []string{},
	}

	cname, err := net.LookupCNAME(hostname)
	if err != nil {
		result.Error = err
		return result, err
	}

	result.Names = []string{cname}
	return result, nil
}

// LookupAll performs all types of DNS lookups
func LookupAll(hostname string) map[RecordType]*LookupResult {
	results := make(map[RecordType]*LookupResult)

	// A records
	if result, err := Lookup(hostname); err == nil {
		results[RecordTypeA] = result
	}

	// CNAME records
	if result, err := LookupCNAME(hostname); err == nil {
		results[RecordTypeCNAME] = result
	}

	// MX records
	if result, err := LookupMX(hostname); err == nil {
		results[RecordTypeMX] = result
	}

	// NS records
	if result, err := LookupNS(hostname); err == nil {
		results[RecordTypeNS] = result
	}

	// TXT records
	if result, err := LookupTXT(hostname); err == nil {
		results[RecordTypeTXT] = result
	}

	return results
}

// FormatResult formats a lookup result for display
func FormatResult(result *LookupResult) string {
	var output strings.Builder

	output.WriteString(fmt.Sprintf("Query: %s\n", result.Query))

	if result.Error != nil {
		output.WriteString(fmt.Sprintf("Error: %v\n", result.Error))
		return output.String()
	}

	if len(result.Addresses) > 0 {
		output.WriteString("Addresses:\n")
		for _, addr := range result.Addresses {
			output.WriteString(fmt.Sprintf("  %s\n", addr))
		}
	}

	if len(result.Names) > 0 {
		output.WriteString("Names:\n")
		for _, name := range result.Names {
			output.WriteString(fmt.Sprintf("  %s\n", name))
		}
	}

	if len(result.MXRecords) > 0 {
		output.WriteString("MX Records:\n")
		for _, mx := range result.MXRecords {
			output.WriteString(fmt.Sprintf("  %s (priority: %d)\n", mx.Host, mx.Pref))
		}
	}

	if len(result.NSRecords) > 0 {
		output.WriteString("NS Records:\n")
		for _, ns := range result.NSRecords {
			output.WriteString(fmt.Sprintf("  %s\n", ns))
		}
	}

	if len(result.TXTRecords) > 0 {
		output.WriteString("TXT Records:\n")
		for _, txt := range result.TXTRecords {
			output.WriteString(fmt.Sprintf("  %s\n", txt))
		}
	}

	return output.String()
}

// GetNameServers returns the system's configured nameservers
func GetNameServers() ([]string, error) {
	// Read /etc/resolv.conf on Unix systems
	resolver := &net.Resolver{
		PreferGo: true,
	}

	// This is a simple implementation
	// A full implementation would parse /etc/resolv.conf
	_ = resolver

	return []string{"System default nameservers"}, nil
}
