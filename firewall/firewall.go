package firewall

import (
	"fmt"
	"net"
	"strings"
	"sync"
)

// Firewall represents a network firewall
type Firewall struct {
	rules []Rule
	mu    sync.RWMutex
}

// Rule represents a firewall rule
type Rule struct {
	ID          int
	Action      Action
	Protocol    string
	SourceIP    string
	SourcePort  uint16
	DestIP      string
	DestPort    uint16
	Description string
}

// Action represents a firewall action
type Action int

const (
	Allow Action = iota
	Deny
	Drop
)

// Packet represents a network packet
type Packet struct {
	Protocol   string
	SourceIP   string
	SourcePort uint16
	DestIP     string
	DestPort   uint16
}

// NewFirewall creates a new firewall
func NewFirewall() *Firewall {
	return &Firewall{
		rules: make([]Rule, 0),
	}
}

// AddRule adds a new firewall rule
func (f *Firewall) AddRule(rule Rule) {
	f.mu.Lock()
	defer f.mu.Unlock()

	rule.ID = len(f.rules) + 1
	f.rules = append(f.rules, rule)
}

// RemoveRule removes a firewall rule by ID
func (f *Firewall) RemoveRule(id int) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	for i, rule := range f.rules {
		if rule.ID == id {
			f.rules = append(f.rules[:i], f.rules[i+1:]...)
			return nil
		}
	}

	return fmt.Errorf("rule %d not found", id)
}

// CheckPacket checks if a packet is allowed by the firewall
func (f *Firewall) CheckPacket(packet Packet) (Action, string) {
	f.mu.RLock()
	defer f.mu.RUnlock()

	// Check each rule in order
	for _, rule := range f.rules {
		if f.matchesRule(packet, rule) {
			return rule.Action, fmt.Sprintf("Matched rule %d: %s", rule.ID, rule.Description)
		}
	}

	// Default action: deny
	return Deny, "No matching rule - default deny"
}

func (f *Firewall) matchesRule(packet Packet, rule Rule) bool {
	// Match protocol
	if rule.Protocol != "" && rule.Protocol != "*" {
		if !strings.EqualFold(packet.Protocol, rule.Protocol) {
			return false
		}
	}

	// Match source IP
	if rule.SourceIP != "" && rule.SourceIP != "*" {
		if !f.matchIP(packet.SourceIP, rule.SourceIP) {
			return false
		}
	}

	// Match destination IP
	if rule.DestIP != "" && rule.DestIP != "*" {
		if !f.matchIP(packet.DestIP, rule.DestIP) {
			return false
		}
	}

	// Match source port
	if rule.SourcePort != 0 && packet.SourcePort != rule.SourcePort {
		return false
	}

	// Match destination port
	if rule.DestPort != 0 && packet.DestPort != rule.DestPort {
		return false
	}

	return true
}

func (f *Firewall) matchIP(packetIP, ruleIP string) bool {
	// Exact match
	if packetIP == ruleIP {
		return true
	}

	// CIDR match
	if strings.Contains(ruleIP, "/") {
		_, network, err := net.ParseCIDR(ruleIP)
		if err != nil {
			return false
		}

		ip := net.ParseIP(packetIP)
		if ip == nil {
			return false
		}

		return network.Contains(ip)
	}

	return false
}

// ListRules returns all firewall rules
func (f *Firewall) ListRules() []Rule {
	f.mu.RLock()
	defer f.mu.RUnlock()

	rules := make([]Rule, len(f.rules))
	copy(rules, f.rules)
	return rules
}

// String returns string representation of action
func (a Action) String() string {
	switch a {
	case Allow:
		return "ALLOW"
	case Deny:
		return "DENY"
	case Drop:
		return "DROP"
	default:
		return "UNKNOWN"
	}
}

// EnableDefaultDeny sets up default deny rules
func (f *Firewall) EnableDefaultDeny() {
	// Allow established connections
	f.AddRule(Rule{
		Action:      Allow,
		Protocol:    "*",
		Description: "Allow established connections",
	})

	// Allow localhost
	f.AddRule(Rule{
		Action:      Allow,
		SourceIP:    "127.0.0.1",
		DestIP:      "127.0.0.1",
		Description: "Allow localhost",
	})
}

// EnableBasicWeb enables basic web traffic
func (f *Firewall) EnableBasicWeb() {
	// Allow HTTP
	f.AddRule(Rule{
		Action:      Allow,
		Protocol:    "tcp",
		DestPort:    80,
		Description: "Allow HTTP",
	})

	// Allow HTTPS
	f.AddRule(Rule{
		Action:      Allow,
		Protocol:    "tcp",
		DestPort:    443,
		Description: "Allow HTTPS",
	})
}
