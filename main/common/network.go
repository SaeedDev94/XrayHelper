package common

import "net"

func IsIPv6(cidr string) bool {
	ip, _, _ := net.ParseCIDR(cidr)
	if ip != nil && ip.To4() == nil {
		return true
	}
	return false
}
