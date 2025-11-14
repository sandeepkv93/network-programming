package arp

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

// Entry represents an ARP table entry
type Entry struct {
	IPAddress  string
	HWAddress  string
	HWType     string
	Flags      string
	Mask       string
	Device     string
}

// GetARPTable retrieves the system's ARP table
func GetARPTable() ([]Entry, error) {
	// Read /proc/net/arp on Linux
	file, err := os.Open("/proc/net/arp")
	if err != nil {
		return nil, fmt.Errorf("failed to open ARP table: %v", err)
	}
	defer file.Close()

	var entries []Entry
	scanner := bufio.NewScanner(file)

	// Skip header line
	scanner.Scan()

	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)

		if len(fields) < 6 {
			continue
		}

		entry := Entry{
			IPAddress: fields[0],
			HWType:    fields[1],
			Flags:     fields[2],
			HWAddress: fields[3],
			Mask:      fields[4],
			Device:    fields[5],
		}

		entries = append(entries, entry)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return entries, nil
}

// GetARPEntry retrieves a specific ARP entry by IP address
func GetARPEntry(ip string) (*Entry, error) {
	entries, err := GetARPTable()
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if entry.IPAddress == ip {
			return &entry, nil
		}
	}

	return nil, fmt.Errorf("ARP entry not found for IP: %s", ip)
}

// FormatEntry formats an ARP entry for display
func FormatEntry(entry Entry) string {
	return fmt.Sprintf("%-15s %-17s %-8s %-6s %-8s %s",
		entry.IPAddress,
		entry.HWAddress,
		entry.HWType,
		entry.Flags,
		entry.Mask,
		entry.Device)
}

// FormatTable formats the entire ARP table for display
func FormatTable(entries []Entry) string {
	output := fmt.Sprintf("%-15s %-17s %-8s %-6s %-8s %s\n",
		"IP Address", "HW Address", "HW Type", "Flags", "Mask", "Device")
	output += strings.Repeat("-", 80) + "\n"

	for _, entry := range entries {
		output += FormatEntry(entry) + "\n"
	}

	return output
}

// ResolveMAC resolves the MAC address for a given IP address
// This is a simpler interface that just returns the MAC address
func ResolveMAC(ip string) (string, error) {
	entry, err := GetARPEntry(ip)
	if err != nil {
		return "", err
	}

	return entry.HWAddress, nil
}

// GetMACByInterface returns the MAC address of a network interface
func GetMACByInterface(interfaceName string) (string, error) {
	iface, err := net.InterfaceByName(interfaceName)
	if err != nil {
		return "", fmt.Errorf("interface not found: %v", err)
	}

	return iface.HardwareAddr.String(), nil
}

// IsComplete checks if an ARP entry is complete (has a valid MAC)
func IsComplete(entry Entry) bool {
	// Check if MAC address is not incomplete (00:00:00:00:00:00)
	return entry.HWAddress != "" &&
		entry.HWAddress != "00:00:00:00:00:00" &&
		!strings.Contains(entry.Flags, "0x0")
}

// GetCompleteEntries returns only complete ARP entries
func GetCompleteEntries() ([]Entry, error) {
	entries, err := GetARPTable()
	if err != nil {
		return nil, err
	}

	var complete []Entry
	for _, entry := range entries {
		if IsComplete(entry) {
			complete = append(complete, entry)
		}
	}

	return complete, nil
}

// GetEntriesByDevice returns ARP entries for a specific network device
func GetEntriesByDevice(device string) ([]Entry, error) {
	entries, err := GetARPTable()
	if err != nil {
		return nil, err
	}

	var filtered []Entry
	for _, entry := range entries {
		if entry.Device == device {
			filtered = append(filtered, entry)
		}
	}

	return filtered, nil
}
