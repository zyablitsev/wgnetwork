package ipcalc

import (
	"encoding/binary"
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

// IP4ToUint32 converts ip bytes to uint32
func IP4ToUint32(ip net.IP) uint32 {
	return binary.BigEndian.Uint32([]byte(ip))
}

// Uint32ToIP4 converts uint32 to ip bytes
func Uint32ToIP4(v uint32) []byte {
	b := make([]byte, 4)
	binary.BigEndian.PutUint32(b, v)
	return net.IP(b)
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

// NextIP returns followed ip address for specified in argument
func NextIP(ip net.IP) net.IP {
	i := IP4ToUint32(ip)
	i++
	return Uint32ToIP4(i)
}

// PrevIP returns previous ip address for specified in argument
func PrevIP(ip net.IP) net.IP {
	i := IP4ToUint32(ip)
	i--
	return Uint32ToIP4(i)
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
