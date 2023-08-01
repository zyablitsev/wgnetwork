package ipcalc

import (
	"encoding/binary"
	"net"
)

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
