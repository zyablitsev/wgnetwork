//go:build linux
// +build linux

package nfutils

import (
	"github.com/google/nftables"
	"github.com/google/nftables/binaryutil"
	"github.com/google/nftables/expr"
)

// SetIIF helper.
func SetIIF(iface string) []expr.Any {
	exprs := []expr.Any{
		ExprLoadIIFName(),
		ExprCmpEqIFName(iface),
	}

	return exprs
}

// SetOIF helper.
func SetOIF(iface string) []expr.Any {
	exprs := []expr.Any{
		ExprLoadOIFName(),
		ExprCmpEqIFName(iface),
	}

	return exprs
}

// SetNIIF helper.
func SetNIIF(iface string) []expr.Any {
	exprs := []expr.Any{
		ExprLoadIIFName(),
		ExprCmpNeqIFName(iface),
	}

	return exprs
}

// SetNOIF helper.
func SetNOIF(iface string) []expr.Any {
	exprs := []expr.Any{
		ExprLoadOIFName(),
		ExprCmpNeqIFName(iface),
	}

	return exprs
}

// SetSourceNet helper.
func SetSourceNet(addr []byte, mask []byte) []expr.Any {
	exprs := []expr.Any{
		ExprLoadNetHeader(1, 12, 4),
		ExprBitwise(1, 1, 4,
			mask,
			[]byte{0, 0, 0, 0}),
		ExprCmpEq(1, addr),
	}

	return exprs
}

// SetProtoICMP helper.
func SetProtoICMP() []expr.Any {
	exprs := []expr.Any{
		ExprLoadNetHeader(1, 9, 1),
		ExprCmpEq(1, ProtoICMP()),
	}

	return exprs
}

// SetICMPTypeEchoRequest helper.
func SetICMPTypeEchoRequest() []expr.Any {
	exprs := []expr.Any{
		ExprLoadTransportHeader(1, 0, 1),
		ExprCmpEq(1, ICMPTypeEchoRequest()),
	}

	return exprs
}

// SetProtoUDP helper.
func SetProtoUDP() []expr.Any {
	exprs := []expr.Any{
		ExprLoadNetHeader(1, 9, 1),
		ExprCmpEq(1, ProtoUDP()),
	}

	return exprs
}

// SetProtoTCP helper.
func SetProtoTCP() []expr.Any {
	exprs := []expr.Any{
		ExprLoadNetHeader(1, 9, 1),
		ExprCmpEq(1, ProtoTCP()),
	}

	return exprs
}

// SetSAddrSet helper.
func SetSAddrSet(s *nftables.Set) []expr.Any {
	exprs := []expr.Any{
		ExprLoadNetHeader(1, 12, 4),
		ExprLookupSet(1, s.Name, s.ID),
	}

	return exprs
}

// SetDAddrSet helper.
func SetDAddrSet(s *nftables.Set) []expr.Any {
	exprs := []expr.Any{
		ExprLoadNetHeader(1, 16, 4),
		ExprLookupSet(1, s.Name, s.ID),
	}

	return exprs
}

// GetAddrSet helper.
func GetAddrSet(t *nftables.Table) *nftables.Set {
	s := &nftables.Set{
		Anonymous: true,
		Constant:  true,
		Table:     t,
		KeyType:   nftables.TypeIPAddr}

	return s
}

// SetSPort helper.
func SetSPort(p uint16) []expr.Any {
	exprs := []expr.Any{
		ExprLoadTransportHeader(1, 0, 2),
		ExprCmpEq(1, binaryutil.BigEndian.PutUint16(p)),
	}

	return exprs
}

// SetDPort helper.
func SetDPort(p uint16) []expr.Any {
	exprs := []expr.Any{
		ExprLoadTransportHeader(1, 2, 2),
		ExprCmpEq(1, binaryutil.BigEndian.PutUint16(p)),
	}

	return exprs
}

// SetSPortSet helper.
func SetSPortSet(s *nftables.Set) []expr.Any {
	exprs := []expr.Any{
		ExprLoadTransportHeader(1, 0, 2),
		ExprLookupSet(1, s.Name, s.ID),
	}

	return exprs
}

// SetDPortSet helper.
func SetDPortSet(s *nftables.Set) []expr.Any {
	exprs := []expr.Any{
		ExprLoadTransportHeader(1, 2, 2),
		ExprLookupSet(1, s.Name, s.ID),
	}

	return exprs
}

// GetPortSet helper.
func GetPortSet(t *nftables.Table) *nftables.Set {
	s := &nftables.Set{
		Anonymous: true,
		Constant:  true,
		Table:     t,
		KeyType:   nftables.TypeInetService}

	return s
}

// GetPortElems helper.
func GetPortElems(ports []uint16) []nftables.SetElement {
	elems := make([]nftables.SetElement, 0, len(ports))
	for _, p := range ports {
		elems = append(elems,
			nftables.SetElement{Key: binaryutil.BigEndian.PutUint16(p)})
	}

	return elems
}

// SetConntrackStateSet helper.
func SetConntrackStateSet(s *nftables.Set) []expr.Any {
	exprs := []expr.Any{
		ExprCtLoadState(1),
		ExprLookupSet(1, s.Name, s.ID),
	}

	return exprs
}

// SetConntrackStateNew helper.
func SetConntrackStateNew() []expr.Any {
	exprs := []expr.Any{
		ExprCtLoadState(1),
		ExprBitwise(1, 1, 4,
			ConntrackStateNew(),
			[]byte{0x00, 0x00, 0x00, 0x00}),
		ExprCmpNeq(1, []byte{0x00, 0x00, 0x00, 0x00}),
	}

	return exprs
}

// SetConntrackStateEstablished helper.
func SetConntrackStateEstablished() []expr.Any {
	exprs := []expr.Any{
		ExprCtLoadState(1),
		ExprBitwise(1, 1, 4,
			ConntrackStateEstablished(),
			[]byte{0x00, 0x00, 0x00, 0x00}),
		ExprCmpNeq(1, []byte{0x00, 0x00, 0x00, 0x00}),
	}

	return exprs
}

// SetConntrackStateRelated helper.
func SetConntrackStateRelated() []expr.Any {
	exprs := []expr.Any{
		ExprCtLoadState(1),
		ExprBitwise(1, 1, 4,
			ConntrackStateRelated(),
			[]byte{0x00, 0x00, 0x00, 0x00}),
		ExprCmpNeq(1, []byte{0x00, 0x00, 0x00, 0x00}),
	}

	return exprs
}

// GetConntrackStateSet helper.
func GetConntrackStateSet(t *nftables.Table) *nftables.Set {
	s := &nftables.Set{
		Anonymous: true,
		Constant:  true,
		Table:     t,
		KeyType:   ConntrackStateDatatype()}

	return s
}

// GetConntrackStateSetElems helper.
func GetConntrackStateSetElems(states []string) []nftables.SetElement {
	elems := make([]nftables.SetElement, 0, len(states))
	for _, s := range states {
		switch s {
		case "new":
			elems = append(elems,
				nftables.SetElement{Key: ConntrackStateNew()})
		case "established":
			elems = append(elems,
				nftables.SetElement{Key: ConntrackStateEstablished()})
		case "related":
			elems = append(elems,
				nftables.SetElement{Key: ConntrackStateRelated()})
		}
	}

	return elems
}
