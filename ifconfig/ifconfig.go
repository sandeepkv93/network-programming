package ifconfig

import (
	"fmt"
	"net"
)

// InterfaceInfo represents information about a network interface
type InterfaceInfo struct {
	Name         string
	HardwareAddr string
	Flags        net.Flags
	Addresses    []AddressInfo
}

// AddressInfo represents an IP address on an interface
type AddressInfo struct {
	IP      string
	Network string
	Netmask string
}

// GetInterfaces retrieves information about all network interfaces
func GetInterfaces() ([]InterfaceInfo, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return nil, fmt.Errorf("failed to get interfaces: %v", err)
	}

	var result []InterfaceInfo

	for _, iface := range interfaces {
		info := InterfaceInfo{
			Name:         iface.Name,
			HardwareAddr: iface.HardwareAddr.String(),
			Flags:        iface.Flags,
			Addresses:    []AddressInfo{},
		}

		// Get addresses for this interface
		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}

		for _, addr := range addrs {
			ipNet, ok := addr.(*net.IPNet)
			if !ok {
				continue
			}

			addressInfo := AddressInfo{
				IP:      ipNet.IP.String(),
				Network: ipNet.Network(),
				Netmask: net.IP(ipNet.Mask).String(),
			}

			info.Addresses = append(info.Addresses, addressInfo)
		}

		result = append(result, info)
	}

	return result, nil
}

// GetInterface retrieves information about a specific interface
func GetInterface(name string) (*InterfaceInfo, error) {
	iface, err := net.InterfaceByName(name)
	if err != nil {
		return nil, fmt.Errorf("interface not found: %v", err)
	}

	info := InterfaceInfo{
		Name:         iface.Name,
		HardwareAddr: iface.HardwareAddr.String(),
		Flags:        iface.Flags,
		Addresses:    []AddressInfo{},
	}

	// Get addresses for this interface
	addrs, err := iface.Addrs()
	if err != nil {
		return &info, nil
	}

	for _, addr := range addrs {
		ipNet, ok := addr.(*net.IPNet)
		if !ok {
			continue
		}

		addressInfo := AddressInfo{
			IP:      ipNet.IP.String(),
			Network: ipNet.Network(),
			Netmask: net.IP(ipNet.Mask).String(),
		}

		info.Addresses = append(info.Addresses, addressInfo)
	}

	return &info, nil
}

// FormatInterface formats interface information for display
func FormatInterface(info InterfaceInfo) string {
	output := fmt.Sprintf("%s: flags=%d<%s>\n", info.Name, info.Flags, info.Flags.String())

	if info.HardwareAddr != "" {
		output += fmt.Sprintf("    ether %s\n", info.HardwareAddr)
	}

	for _, addr := range info.Addresses {
		if addr.Network == "ip+net" {
			output += fmt.Sprintf("    inet %s netmask %s\n", addr.IP, addr.Netmask)
		}
	}

	return output
}

// IsUp checks if an interface is up
func IsUp(name string) (bool, error) {
	iface, err := net.InterfaceByName(name)
	if err != nil {
		return false, err
	}

	return iface.Flags&net.FlagUp != 0, nil
}

// IsLoopback checks if an interface is a loopback interface
func IsLoopback(name string) (bool, error) {
	iface, err := net.InterfaceByName(name)
	if err != nil {
		return false, err
	}

	return iface.Flags&net.FlagLoopback != 0, nil
}
