package ipcalc

import (
	"fmt"
	"math"
	"net"
)

// ParseCIDR is like net.ParseCIDR except that it parses IPv4 addresses as 4
// byte addresses instead of 16-byte mapped IPv6 addresses. Much like ParseIP.
func ParseCIDR(cidr string) (net.IP, *net.IPNet, error) {
	ip, n, err := net.ParseCIDR(cidr)
	if err != nil {
		return nil, nil, err
	}

	ip = ip.To4()
	if ip == nil {
		return nil, nil, nil
	}

	return ip, &net.IPNet{IP: n.IP.To4(), Mask: n.Mask}, nil
}

// FirstIP returns the network address.
func FirstIP(n *net.IPNet) net.IP {
	if len(n.IP) != 4 {
		return nil
	}

	ip := net.IPv4(0, 0, 0, 0).To4()
	for i := 0; i < len(n.IP); i++ {
		ip[i] = n.IP[i] & n.Mask[i]
	}

	return ip
}

// LastIP returns the broadcast address.
func LastIP(n *net.IPNet) net.IP {
	if len(n.IP) != 4 {
		return nil
	}

	ip := net.IPv4(0, 0, 0, 0).To4()
	for i := 0; i < len(n.IP); i++ {
		ip[i] = n.IP[i] | ^n.Mask[i]
	}

	return ip
}

// FirstHostIP returns first host address for the given network.
func FirstHostIP(n *net.IPNet) net.IP {
	ip := FirstIP(n)
	maskLen, bits := n.Mask.Size()
	hostLen := uint32(bits) - uint32(maskLen)
	if hostLen > 1 {
		ip = NextIP(ip)
	}

	return ip
}

// LastHostIP returns last host address for the given network.
func LastHostIP(n *net.IPNet) net.IP {
	ip := LastIP(n)
	maskLen, bits := n.Mask.Size()
	hostLen := uint32(bits) - uint32(maskLen)
	if hostLen > 1 {
		ip = PrevIP(ip)
	}

	return ip
}

// NetworkSize returns amount of all network addresses as first value
// and amount of host addresses as second value.
func NetworkSize(n *net.IPNet) (uint32, uint32) {
	maskLen, bits := n.Mask.Size()
	hostLen := uint32(bits) - uint32(maskLen)
	total := uint32(1) << hostLen
	hosts := total
	if hostLen > 1 {
		hosts -= 2
	}

	return total, hosts
}

// MinMaxAddr returns first and last network addresses.
func MinMaxAddr(network *net.IPNet) (net.IP, net.IP) {
	maskLen, bits := network.Mask.Size()
	if bits != 32 {
		return nil, nil
	}

	firstip := FirstIP(network)
	if maskLen == bits {
		min := firstip
		max := make([]byte, len(min))
		copy(max, min)
		return min, max
	}

	firstipint := IP4ToUint32(firstip)
	hostLen := uint32(bits) - uint32(maskLen)
	lastipint := ((uint32(1) << hostLen) - uint32(1)) | firstipint
	lastip := Uint32ToIP4(lastipint)

	return firstip, lastip
}

// MinMaxHost returns first and last host addresses.
func MinMaxHost(network *net.IPNet) (net.IP, net.IP) {
	minAddr, maxAddr := MinMaxAddr(network)
	addrsamount, _ := NetworkSize(network)
	if addrsamount < 3 {
		return minAddr, maxAddr
	}

	return NextIP(minAddr), PrevIP(maxAddr)
}

// NetworkIPs returns full list of ip addresses for specified network.
// Including network and broadcast ip addresses.
func NetworkIPs(n *net.IPNet, limit uint32) []net.IP {
	size, _ := NetworkSize(n)
	if limit > 0 && limit < size {
		size = limit
	}

	result := make([]net.IP, size)
	next := FirstIP(n)
	for i := uint32(0); i < size; i++ {
		result[i] = next[:]
		next = NextIP(next)
	}

	return result
}

// NetworkHostIPs returns full list of hosts ip addresses for specified network.
// Excluding network and broadcast ip addresses.
func NetworkHostIPs(n *net.IPNet, limit uint32) []net.IP {
	_, size := NetworkSize(n)
	if limit > 0 && limit < size {
		size = limit
	}

	result := make([]net.IP, size)
	next := FirstHostIP(n)
	for i := uint32(0); i < size; i++ {
		result[i] = next[:]
		next = NextIP(next)
	}

	return result
}

// IsNetIntersected checks two networks for intersection.
func IsNetIntersected(n1, n2 net.IPNet) bool {
	return n2.Contains(n1.IP) || n1.Contains(n2.IP)
}

// IP4RangeToCIDRs converts ip v4 range to list of cidr's
func IP4RangeToCIDRs(startIP, endIP net.IP) ([]net.IPNet, error) {
	startUint32IP := IP4ToUint32(startIP)
	endUint32IP := IP4ToUint32(endIP)

	if startUint32IP > endUint32IP {
		err := fmt.Errorf(
			"start ip %q must be lte end ip %q",
			startIP, endIP)

		return nil, err
	}

	cidr2mask := []uint32{
		0x00000000, 0x80000000, 0xC0000000,
		0xE0000000, 0xF0000000, 0xF8000000,
		0xFC000000, 0xFE000000, 0xFF000000,
		0xFF800000, 0xFFC00000, 0xFFE00000,
		0xFFF00000, 0xFFF80000, 0xFFFC0000,
		0xFFFE0000, 0xFFFF0000, 0xFFFF8000,
		0xFFFFC000, 0xFFFFE000, 0xFFFFF000,
		0xFFFFF800, 0xFFFFFC00, 0xFFFFFE00,
		0xFFFFFF00, 0xFFFFFF80, 0xFFFFFFC0,
		0xFFFFFFE0, 0xFFFFFFF0, 0xFFFFFFF8,
		0xFFFFFFFC, 0xFFFFFFFE, 0xFFFFFFFF,
	}

	cidrs := []net.IPNet{}
	for endUint32IP >= startUint32IP {
		maxSize := 32
		for maxSize > 0 {
			maskedBase := startUint32IP & cidr2mask[maxSize-1]
			if maskedBase != startUint32IP {
				break
			}

			maxSize--
		}

		x := math.Log(float64(endUint32IP-startUint32IP+1)) / math.Log(2)
		maxDiff := 32 - int(math.Floor(x))
		if maxSize < maxDiff {
			maxSize = maxDiff
		}

		cidrs = append(
			cidrs,
			net.IPNet{
				IP:   Uint32ToIP4(startUint32IP),
				Mask: net.CIDRMask(maxSize, 32),
			},
		)

		startUint32IP += uint32(math.Exp2(float64(32 - maxSize)))
	}

	return cidrs, nil
}
