//go:build linux
// +build linux

package nfutils

import "github.com/google/nftables"

// ProtoICMP bytes.
func ProtoICMP() []byte {
	return []byte{0x01}
}

// ICMPTypeEchoRequest bytes.
func ICMPTypeEchoRequest() []byte {
	return []byte{0x08}
}

// ProtoUDP bytes.
func ProtoUDP() []byte {
	return []byte{0x11}
}

// ProtoTCP bytes.
func ProtoTCP() []byte {
	return []byte{0x06}
}

// ConntrackStateNew bytes.
func ConntrackStateNew() []byte {
	return []byte{0x08, 0x00, 0x00, 0x00}
}

// ConntrackStateEstablished bytes.
func ConntrackStateEstablished() []byte {
	return []byte{0x02, 0x00, 0x00, 0x00}
}

// ConntrackStateRelated bytes.
func ConntrackStateRelated() []byte {
	return []byte{0x04, 0x00, 0x00, 0x00}
}

// ConntrackStateDatatype object.
func ConntrackStateDatatype() nftables.SetDatatype {
	ctStateDataType := nftables.SetDatatype{Name: "ct_state", Bytes: 4}
	// nftMagic: https://git.netfilter.org/nftables/tree/src/datatype.c#n32 (arr index)
	ctStateDataType.SetNFTMagic(26)
	return ctStateDataType
}
