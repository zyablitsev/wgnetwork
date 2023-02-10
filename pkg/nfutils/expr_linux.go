// +build linux

package nfutils

import (
	"net"

	"golang.org/x/sys/unix"

	"github.com/google/nftables/expr"
)

// ExprLoadIIFName wrapper
func ExprLoadIIFName() *expr.Meta {
	// [ meta load iifname => reg 1 ]
	return &expr.Meta{Key: expr.MetaKeyIIFNAME, Register: 1}
}

// ExprLoadOIFName wrapper
func ExprLoadOIFName() *expr.Meta {
	// [ meta load oifname => reg 1 ]
	return &expr.Meta{Key: expr.MetaKeyOIFNAME, Register: 1}
}

// ExprCmpEqIFName wrapper
func ExprCmpEqIFName(name string) *expr.Cmp {
	// [ cmp eq reg 1 0x00006f6c 0x00000000 0x00000000 0x00000000 ]
	return &expr.Cmp{
		Op:       expr.CmpOpEq,
		Register: 1,
		Data:     ifname(name),
	}
}

// ExprCmpNeqIFName wrapper
func ExprCmpNeqIFName(name string) *expr.Cmp {
	// [ cmp neq reg 1 0x00006f6c 0x00000000 0x00000000 0x00000000 ]
	return &expr.Cmp{
		Op:       expr.CmpOpNeq,
		Register: 1,
		Data:     ifname(name),
	}
}

func ifname(n string) []byte {
	b := make([]byte, 16)
	copy(b, []byte(n+"\x00"))
	return b
}

// ExprLoadNetHeader wrapper
func ExprLoadNetHeader(reg, offset, l uint32) *expr.Payload {
	// [ payload load 4b @ network header + 12 => reg 1 ]
	return &expr.Payload{
		DestRegister: reg,
		Base:         expr.PayloadBaseNetworkHeader,
		Offset:       offset,
		Len:          l,
	}
}

// ExprLoadTransportHeader wrapper
func ExprLoadTransportHeader(reg, offset, l uint32) *expr.Payload {
	// [ payload load 1b @ transport header + 0 => reg 1 ]
	return &expr.Payload{
		DestRegister: reg,
		Base:         expr.PayloadBaseTransportHeader,
		Offset:       offset,
		Len:          l,
	}
}

// ExprBitwise wrapper
func ExprBitwise(dReg, sReg, l uint32, mask, xor []byte) *expr.Bitwise {
	// [ bitwise reg 1 = (reg=1 & 0x000000ff ) ^ 0x00000000 ]
	return &expr.Bitwise{
		DestRegister:   dReg,
		SourceRegister: sReg,
		Len:            l,
		Mask:           mask,
		Xor:            xor,
	}
}

// ExprCmpEq wrapper
func ExprCmpEq(reg uint32, data []byte) *expr.Cmp {
	// [ cmp eq reg 1 0x0000007f ]
	return &expr.Cmp{
		Op:       expr.CmpOpEq,
		Register: reg,
		Data:     data,
	}
}

// ExprCmpNeq wrapper
func ExprCmpNeq(reg uint32, data []byte) *expr.Cmp {
	// [ cmp eq reg 1 0x0000007f ]
	return &expr.Cmp{
		Op:       expr.CmpOpNeq,
		Register: reg,
		Data:     data,
	}
}

// ExprLookupSet wrapper
func ExprLookupSet(reg uint32, name string, id uint32) *expr.Lookup {
	// [ lookup reg 1 set adminipset ]
	return &expr.Lookup{
		SourceRegister: 1,
		SetName:        name,
		SetID:          id,
	}
}

// ExprCtLoadState wrapper
func ExprCtLoadState(reg uint32) *expr.Ct {
	// [ ct load state => reg 1 ]
	return &expr.Ct{
		// Key:      unix.NFT_CT_STATE,
		Key:      expr.CtKeySTATE,
		Register: reg,
	}
}

// ExprImmediate wrapper
func ExprImmediate(ip net.IP) *expr.Immediate {
	// [ immediate reg 1 0x0158a8c0 ]
	return &expr.Immediate{
		Register: 1,
		Data:     ip,
	}
}

// ExprSNAT wrapper
func ExprSNAT(addrMin, addrMax uint32) *expr.NAT {
	// [ nat snat ip addr_min reg 1 addr_max reg 0 ]
	return &expr.NAT{
		Type:       expr.NATTypeSourceNAT,
		Family:     unix.NFPROTO_IPV4,
		RegAddrMin: addrMin,
		RegAddrMax: addrMax,
	}
}

// ExprAccept wrapper
func ExprAccept() *expr.Verdict {
	// [ immediate reg 0 accept ]
	return &expr.Verdict{
		Kind: expr.VerdictAccept,
	}
}

// ExprDrop wrapper
func ExprDrop() *expr.Verdict {
	// [ immediate reg 0 accept ]
	return &expr.Verdict{
		Kind: expr.VerdictDrop,
	}
}

// ExprReject wrapper
func ExprReject(t uint32, c uint8) *expr.Reject {
	// [ reject type 0 code 3 ]
	return &expr.Reject{Type: t, Code: c}
}
