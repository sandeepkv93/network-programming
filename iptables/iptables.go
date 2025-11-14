package iptables

import (
	"fmt"
	"strings"
	"sync"
)

// IPTables represents an iptables-like firewall system
type IPTables struct {
	chains map[string]*Chain
	mu     sync.RWMutex
}

// Chain represents a chain of rules (INPUT, OUTPUT, FORWARD)
type Chain struct {
	Name   string
	Rules  []Rule
	Policy Policy
	mu     sync.Mutex
}

// Rule represents an iptables rule
type Rule struct {
	Table      string // filter, nat, mangle
	Chain      string // INPUT, OUTPUT, FORWARD, PREROUTING, POSTROUTING
	Protocol   string // tcp, udp, icmp, all
	Source     string // source IP/CIDR
	Dest       string // destination IP/CIDR
	Sport      string // source port
	Dport      string // destination port
	Target     string // ACCEPT, DROP, REJECT, MASQUERADE
	Interface  string // input/output interface
	State      string // NEW, ESTABLISHED, RELATED
}

// Policy represents default chain policy
type Policy string

const (
	ACCEPT Policy = "ACCEPT"
	DROP   Policy = "DROP"
	REJECT Policy = "REJECT"
)

// NewIPTables creates a new IPTables instance
func NewIPTables() *IPTables {
	ipt := &IPTables{
		chains: make(map[string]*Chain),
	}

	// Initialize default chains
	ipt.chains["INPUT"] = &Chain{Name: "INPUT", Policy: ACCEPT, Rules: []Rule{}}
	ipt.chains["OUTPUT"] = &Chain{Name: "OUTPUT", Policy: ACCEPT, Rules: []Rule{}}
	ipt.chains["FORWARD"] = &Chain{Name: "FORWARD", Policy: DROP, Rules: []Rule{}}

	return ipt
}

// AddRule adds a rule to a chain
func (ipt *IPTables) AddRule(rule Rule) error {
	ipt.mu.Lock()
	defer ipt.mu.Unlock()

	chain, exists := ipt.chains[rule.Chain]
	if !exists {
		return fmt.Errorf("chain %s does not exist", rule.Chain)
	}

	chain.mu.Lock()
	defer chain.mu.Unlock()

	chain.Rules = append(chain.Rules, rule)
	return nil
}

// DeleteRule deletes a rule from a chain by index
func (ipt *IPTables) DeleteRule(chainName string, index int) error {
	ipt.mu.Lock()
	defer ipt.mu.Unlock()

	chain, exists := ipt.chains[chainName]
	if !exists {
		return fmt.Errorf("chain %s does not exist", chainName)
	}

	chain.mu.Lock()
	defer chain.mu.Unlock()

	if index < 0 || index >= len(chain.Rules) {
		return fmt.Errorf("invalid rule index")
	}

	chain.Rules = append(chain.Rules[:index], chain.Rules[index+1:]...)
	return nil
}

// SetPolicy sets the default policy for a chain
func (ipt *IPTables) SetPolicy(chainName string, policy Policy) error {
	ipt.mu.Lock()
	defer ipt.mu.Unlock()

	chain, exists := ipt.chains[chainName]
	if !exists {
		return fmt.Errorf("chain %s does not exist", chainName)
	}

	chain.mu.Lock()
	defer chain.mu.Unlock()

	chain.Policy = policy
	return nil
}

// ListRules lists all rules in a chain
func (ipt *IPTables) ListRules(chainName string) ([]Rule, error) {
	ipt.mu.RLock()
	defer ipt.mu.RUnlock()

	chain, exists := ipt.chains[chainName]
	if !exists {
		return nil, fmt.Errorf("chain %s does not exist", chainName)
	}

	chain.mu.Lock()
	defer chain.mu.Unlock()

	rules := make([]Rule, len(chain.Rules))
	copy(rules, chain.Rules)
	return rules, nil
}

// FlushChain removes all rules from a chain
func (ipt *IPTables) FlushChain(chainName string) error {
	ipt.mu.Lock()
	defer ipt.mu.Unlock()

	chain, exists := ipt.chains[chainName]
	if !exists {
		return fmt.Errorf("chain %s does not exist", chainName)
	}

	chain.mu.Lock()
	defer chain.mu.Unlock()

	chain.Rules = []Rule{}
	return nil
}

// String returns a string representation of a rule (iptables-like syntax)
func (r Rule) String() string {
	var parts []string

	parts = append(parts, "-A", r.Chain)

	if r.Protocol != "" && r.Protocol != "all" {
		parts = append(parts, "-p", r.Protocol)
	}

	if r.Source != "" {
		parts = append(parts, "-s", r.Source)
	}

	if r.Dest != "" {
		parts = append(parts, "-d", r.Dest)
	}

	if r.Sport != "" {
		parts = append(parts, "--sport", r.Sport)
	}

	if r.Dport != "" {
		parts = append(parts, "--dport", r.Dport)
	}

	if r.Interface != "" {
		parts = append(parts, "-i", r.Interface)
	}

	if r.State != "" {
		parts = append(parts, "-m", "state", "--state", r.State)
	}

	if r.Target != "" {
		parts = append(parts, "-j", r.Target)
	}

	return strings.Join(parts, " ")
}

// Example presets

// AllowSSH adds a rule to allow SSH
func (ipt *IPTables) AllowSSH() {
	ipt.AddRule(Rule{
		Chain:    "INPUT",
		Protocol: "tcp",
		Dport:    "22",
		Target:   "ACCEPT",
	})
}

// AllowHTTP adds rules to allow HTTP and HTTPS
func (ipt *IPTables) AllowHTTP() {
	ipt.AddRule(Rule{
		Chain:    "INPUT",
		Protocol: "tcp",
		Dport:    "80",
		Target:   "ACCEPT",
	})
	ipt.AddRule(Rule{
		Chain:    "INPUT",
		Protocol: "tcp",
		Dport:    "443",
		Target:   "ACCEPT",
	})
}

// AllowEstablished allows established connections
func (ipt *IPTables) AllowEstablished() {
	ipt.AddRule(Rule{
		Chain:  "INPUT",
		State:  "ESTABLISHED,RELATED",
		Target: "ACCEPT",
	})
}
