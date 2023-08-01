//go:build linux
// +build linux

package nfutils

import (
	"bytes"
	"errors"
	"math"
	"net"

	"github.com/google/nftables"
	"github.com/google/nftables/binaryutil"
	"github.com/google/nftables/expr"
	"golang.org/x/sys/unix"

	"wgnetwork/pkg/ipcalc"
)

// IFaceExpressions build list of nftables exprs for ifaces.
func IFaceExpressions(
	c *nftables.Conn, t *nftables.Table,
	metaKey expr.MetaKey,
	ifnames []string, invert bool,
) ([]expr.Any, error) {
	if len(ifnames) == 0 {
		return nil, nil
	}

	if len(ifnames) == 1 {
		var op expr.CmpOp = expr.CmpOpEq
		if invert {
			op = expr.CmpOpNeq
		}

		exprs := []expr.Any{
			&expr.Meta{Key: metaKey, Register: 1},
			&expr.Cmp{Op: op, Register: 1, Data: ifname(ifnames[0])},
		}

		return exprs, nil
	}

	items := make([]nftables.SetElement, len(ifnames))
	for i, name := range ifnames {
		items[i] = nftables.SetElement{Key: ifname(name)}
	}

	s := &nftables.Set{
		Anonymous: true,
		Constant:  true,
		Table:     t,
		KeyType:   nftables.TypeIFName}

	err := c.AddSet(s, items)
	if err != nil {
		return nil, err
	}

	exprs := []expr.Any{
		&expr.Meta{Key: metaKey, Register: 1},
		&expr.Lookup{
			SourceRegister: 1,
			SetName:        s.Name,
			SetID:          s.ID,
			Invert:         invert,
		},
	}

	return exprs, nil
}

func ifname(n string) []byte {
	b := make([]byte, 16)
	copy(b, []byte(n+"\x00"))
	return b
}

// L4ProtoExpressions build list of nftables exprs for protos.
func L4ProtoExpressions(
	c *nftables.Conn, t *nftables.Table,
	protoNames []string, invert bool,
) ([]expr.Any, error) {
	if len(protoNames) == 0 {
		return nil, nil
	}

	if len(protoNames) == 1 {
		key, err := l4protoMagic(protoNames[0])
		if err != nil {
			return nil, err
		}

		var op expr.CmpOp = expr.CmpOpEq
		if invert {
			op = expr.CmpOpNeq
		}

		exprs := []expr.Any{
			// [ meta load l4proto => reg 1 ]
			&expr.Meta{Key: expr.MetaKeyL4PROTO, Register: 1},
			// [ cmp eq reg 1 0x00000006 ]
			&expr.Cmp{Op: op, Register: 1, Data: key},
		}

		return exprs, nil
	}

	items := make([]nftables.SetElement, len(protoNames))
	for i, name := range protoNames {
		key, err := l4protoMagic(name)
		if err != nil {
			return nil, err
		}

		items[i] = nftables.SetElement{Key: key}
	}

	s := &nftables.Set{
		Anonymous: true,
		Constant:  true,
		Table:     t,
		KeyType:   nftables.TypeInetProto}

	err := c.AddSet(s, items)
	if err != nil {
		return nil, err
	}

	exprs := []expr.Any{
		&expr.Meta{Key: expr.MetaKeyL4PROTO, Register: 1},
		&expr.Lookup{
			SourceRegister: 1,
			SetName:        s.Name,
			SetID:          s.ID,
			Invert:         invert,
		},
	}

	return exprs, nil
}

func l4protoMagic(name string) ([]byte, error) {
	var key []byte
	switch name {
	case "udp":
		key = []byte{unix.IPPROTO_UDP}
	case "tcp":
		key = []byte{unix.IPPROTO_TCP}
	case "icmp":
		key = []byte{unix.IPPROTO_ICMP}
	default:
		return nil, errors.New("unsupported proto")
	}

	return key, nil
}

// ICMPTypeExpressions build list of nftables exprs for icmp type.
func ICMPTypeExpressions(
	c *nftables.Conn, t *nftables.Table,
	icmpType []byte, invert bool,
) ([]expr.Any, error) {
	if icmpType == nil {
		return nil, nil
	}

	var op expr.CmpOp = expr.CmpOpEq
	if invert {
		op = expr.CmpOpNeq
	}

	exprs := []expr.Any{
		&expr.Payload{
			DestRegister: 1,
			Base:         expr.PayloadBaseTransportHeader,
			Offset:       0,
			Len:          1,
		},
		&expr.Cmp{
			Op:       op,
			Register: 1,
			Data:     icmpType,
		},
	}

	return exprs, nil
}

// TransportExpressions build list of nftables exprs for ports.
func TransportExpressions(
	c *nftables.Conn, t *nftables.Table,
	ports [][2]uint16, invert bool, isDest bool,
) ([]expr.Any, error) {
	if len(ports) == 0 {
		return nil, nil
	}

	offset := uint32(0)
	if isDest {
		offset = 2
	}

	if len(ports) == 1 {
		if ports[0][1] == 0 {
			var op expr.CmpOp = expr.CmpOpEq
			if invert {
				op = expr.CmpOpNeq
			}

			key := binaryutil.BigEndian.PutUint16(ports[0][0])
			exprs := []expr.Any{
				&expr.Payload{
					DestRegister: 1,
					Base:         expr.PayloadBaseTransportHeader,
					Offset:       offset,
					Len:          2,
				},
				&expr.Cmp{Op: op, Register: 1, Data: key},
			}

			return exprs, nil
		}

		gte := binaryutil.BigEndian.PutUint16(ports[0][0])
		lte := binaryutil.BigEndian.PutUint16(ports[0][1])

		if invert {
			exprs := []expr.Any{
				&expr.Payload{
					DestRegister: 1,
					Base:         expr.PayloadBaseTransportHeader,
					Offset:       offset,
					Len:          2,
				},
				&expr.Range{
					Op:       expr.CmpOpNeq,
					Register: 1,
					FromData: gte,
					ToData:   lte,
				},
			}

			return exprs, nil
		}

		exprs := []expr.Any{
			&expr.Payload{
				DestRegister: 1,
				Base:         expr.PayloadBaseTransportHeader,
				Offset:       offset,
				Len:          2,
			},
			&expr.Cmp{Op: expr.CmpOpGte, Register: 1, Data: gte},
			&expr.Cmp{Op: expr.CmpOpLte, Register: 1, Data: lte},
		}

		return exprs, nil
	}

	var isInterval bool = false
	for _, port := range ports {
		if port[1] > port[0] {
			isInterval = true
		}
	}

	cnt := len(ports)
	if isInterval {
		cnt *= 2
	}

	items := make([]nftables.SetElement, 0, cnt)
	for _, port := range ports {
		if port[1] == 0 {
			if isInterval {
				items = append(
					items,
					nftables.SetElement{
						Key: binaryutil.BigEndian.PutUint16(port[0]),
					},
				)
				if port[0] < math.MaxUint16 {
					items = append(
						items,
						nftables.SetElement{
							Key:         binaryutil.BigEndian.PutUint16(port[0] + 1),
							IntervalEnd: true,
						},
					)
				}
			} else {
				items = append(
					items,
					nftables.SetElement{
						Key: binaryutil.BigEndian.PutUint16(port[0]),
					},
				)
			}

			continue
		}

		items = append(
			items,
			nftables.SetElement{
				Key: binaryutil.BigEndian.PutUint16(port[0]),
			},
		)
		if port[1] < math.MaxUint16 {
			items = append(
				items,
				nftables.SetElement{
					Key:         binaryutil.BigEndian.PutUint16(port[1] + 1),
					IntervalEnd: true,
				},
			)
		}
	}

	s := &nftables.Set{
		Anonymous:     true,
		Constant:      true,
		Table:         t,
		KeyType:       nftables.TypeInetService,
		Interval:      isInterval,
		Concatenation: false,
	}

	err := c.AddSet(s, items)
	if err != nil {
		return nil, err
	}

	exprs := []expr.Any{
		&expr.Payload{
			DestRegister: 1,
			Base:         expr.PayloadBaseTransportHeader,
			Offset:       offset,
			Len:          2,
		},
		&expr.Lookup{
			SourceRegister: 1,
			SetName:        s.Name,
			SetID:          s.ID,
			Invert:         invert,
		},
	}

	return exprs, nil
}

// NetworkExpressions build list of nftables exprs for networks.
func NetworkExpressions(
	c *nftables.Conn, t *nftables.Table,
	namedSet *nftables.Set, ipnets [][2]*net.IPNet, invert bool, isDest bool,
) ([]expr.Any, error) {
	offset := uint32(12)
	if isDest {
		offset = 16
	}

	if namedSet != nil {
		exprs := []expr.Any{
			&expr.Payload{
				DestRegister: 1,
				Base:         expr.PayloadBaseNetworkHeader,
				Offset:       offset,
				Len:          4,
			},
			&expr.Lookup{
				SourceRegister: 1,
				SetName:        namedSet.Name,
				SetID:          namedSet.ID,
				Invert:         invert,
			},
		}

		return exprs, nil
	}

	if len(ipnets) == 0 {
		return nil, nil
	}

	if len(ipnets) == 1 {
		if ipnets[0][1] == nil {
			var op expr.CmpOp = expr.CmpOpEq
			if invert {
				op = expr.CmpOpNeq
			}

			key := ipnets[0][0].IP.To4()
			mask := ipnets[0][0].Mask
			exprs := []expr.Any{
				&expr.Payload{
					DestRegister: 1,
					Base:         expr.PayloadBaseNetworkHeader,
					Offset:       offset,
					Len:          4,
				},
				&expr.Bitwise{
					DestRegister:   1,
					SourceRegister: 1,
					Len:            4,
					Mask:           mask,
					Xor:            []byte{0, 0, 0, 0},
				},
				&expr.Cmp{Op: op, Register: 1, Data: key},
			}

			return exprs, nil
		}

		gte := ipnets[0][0].IP.To4()
		lte := ipnets[0][1].IP.To4()
		// lte := ipcalc.NextIP(ipnets[0][1].To4()).To4()

		if invert {
			exprs := []expr.Any{
				&expr.Payload{
					DestRegister: 1,
					Base:         expr.PayloadBaseNetworkHeader,
					Offset:       offset,
					Len:          4,
				},
				&expr.Range{
					Op:       expr.CmpOpNeq,
					Register: 1,
					FromData: gte,
					ToData:   lte,
				},
			}

			return exprs, nil
		}

		exprs := []expr.Any{
			&expr.Payload{
				DestRegister: 1,
				Base:         expr.PayloadBaseNetworkHeader,
				Offset:       offset,
				Len:          4,
			},
			&expr.Cmp{Op: expr.CmpOpGte, Register: 1, Data: gte},
			&expr.Cmp{Op: expr.CmpOpLte, Register: 1, Data: lte},
		}

		return exprs, nil
	}

	var isInterval bool = false
	for _, ipnet := range ipnets {
		if ipnet[1] != nil ||
			!bytes.Equal(ipnet[0].Mask, net.IPv4Mask(255, 255, 255, 255)) {

			isInterval = true
		}
	}

	cnt := len(ipnets)
	if isInterval {
		cnt *= 2
	}

	items := make([]nftables.SetElement, 0, cnt)
	for _, ipnet := range ipnets {
		if isInterval {
			ip := ipnet[0].IP.To4()
			mask := ipnet[0].Mask
			items = append(
				items,
				nftables.SetElement{Key: ip},
			)

			if ipnet[1] != nil {
				next := ipcalc.NextIP(ipnet[1].IP.To4()).To4()
				items = append(
					items,
					nftables.SetElement{
						Key:         next,
						IntervalEnd: true,
					},
				)
			} else if bytes.Equal(mask, net.IPv4Mask(255, 255, 255, 255)) {
				next := ipcalc.NextIP(ipnet[0].IP.To4()).To4()
				items = append(
					items,
					nftables.SetElement{
						Key:         next,
						IntervalEnd: true,
					},
				)
			} else {
				_, maxip := ipcalc.MinMaxAddr(ipnet[0])
				next := ipcalc.NextIP(maxip.To4()).To4()
				items = append(
					items,
					nftables.SetElement{
						Key:         next,
						IntervalEnd: true,
					},
				)
			}
		} else {
			items = append(
				items,
				nftables.SetElement{Key: ipnet[0].IP.To4()},
			)
		}
	}

	s := &nftables.Set{
		Anonymous:     true,
		Constant:      true,
		Table:         t,
		KeyType:       nftables.TypeIPAddr,
		Interval:      isInterval,
		Concatenation: false,
	}

	err := c.AddSet(s, items)
	if err != nil {
		return nil, err
	}

	exprs := []expr.Any{
		&expr.Payload{
			DestRegister: 1,
			Base:         expr.PayloadBaseNetworkHeader,
			Offset:       offset,
			Len:          4,
		},
		&expr.Lookup{
			SourceRegister: 1,
			SetName:        s.Name,
			SetID:          s.ID,
			Invert:         invert,
		},
	}

	return exprs, nil
}

// ConntrackExpressions build list of nftables exprs for connection tracking states.
func ConntrackExpressions(
	c *nftables.Conn, t *nftables.Table,
	states []string,
) ([]expr.Any, error) {
	if len(states) == 0 {
		return nil, nil
	}

	var state uint32 = 0
	for _, s := range states {
		switch s {
		case "untracked":
			state |= expr.CtStateBitUNTRACKED
		case "new":
			state |= expr.CtStateBitNEW
		case "established":
			state |= expr.CtStateBitESTABLISHED
		case "related":
			state |= expr.CtStateBitRELATED
		case "invalid":
			state |= expr.CtStateBitINVALID
		default:
			return nil, errors.New("unsupported conntrack state")
		}
	}

	exprs := []expr.Any{
		&expr.Ct{Register: 1, SourceRegister: false, Key: expr.CtKeySTATE},
		&expr.Bitwise{
			SourceRegister: 1,
			DestRegister:   1,
			Len:            4,
			Mask:           binaryutil.NativeEndian.PutUint32(state),
			Xor:            binaryutil.NativeEndian.PutUint32(0),
		},
		&expr.Cmp{Op: expr.CmpOpNeq, Register: 1, Data: []byte{0, 0, 0, 0}},
	}

	return exprs, nil
}

// AcceptExpression builder
func AcceptExpression() *expr.Verdict {
	return &expr.Verdict{
		Kind: expr.VerdictAccept,
	}
}

// DropExpression builder
func DropExpression() *expr.Verdict {
	return &expr.Verdict{
		Kind: expr.VerdictDrop,
	}
}

// JumpExpression builder
func JumpExpression(chain string) *expr.Verdict {
	return &expr.Verdict{
		Kind:  expr.VerdictJump,
		Chain: chain,
	}
}

// RejectExpression builder
func RejectExpression(t uint32, c uint8) *expr.Reject {
	return &expr.Reject{Type: t, Code: c}
}
