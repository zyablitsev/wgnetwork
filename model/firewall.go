package model

import (
	"net"

	"github.com/google/nftables"
)

// FilterRule nftables.
type FilterRule struct {
	IIF    []string `json:"iiface"`
	IIFNeq bool     `json:"iiface_neq"`

	OIF    []string `json:"oiface"`
	OIFNeq bool     `json:"oiface_neq"`

	Proto    []string `json:"proto"`
	ProtoNeq bool     `json:"proto_neq"`

	ICMPType    []byte `json:"icmp_type"`
	ICMPTypeNeq bool   `json:"icmp_type_neq"`

	SAddr    [][2]*net.IPNet `json:"saddr"`
	SAddrSet *nftables.Set   `json:"saddr_set"`
	SAddrNeq bool            `json:"saddr_neq"`

	SPort    [][2]uint16 `json:"sport"`
	SPortNeq bool        `json:"sport_neq"`

	DAddr    [][2]*net.IPNet `json:"daddr"`
	DAddrSet *nftables.Set   `json:"daddr_set"`
	DAddrNeq bool            `json:"daddr_neq"`

	DPort    [][2]uint16 `json:"dport"`
	DPortNeq bool        `json:"dport_neq"`

	Conntrack []string `json:"conntrack"`

	Action           string `json:"action"`
	JumpActionChain  string `json:"jump_action_chain"`
	RejectActionType uint32 `json:"reject_action_type"`
	RejectActionCode uint8  `json:"reject_action_code"`
}

// SNATRule nftables.
type SNATRule struct {
	IIF []string `json:"iiface"`
	OIF []string `json:"oiface"`

	SAddr    [][2]*net.IPNet `json:"saddr"`
	SAddrSet *nftables.Set   `json:"saddr_set"`

	SNATToIP net.IP `json:"snat_to_ip"`
}

// DNATRule nftables.
type DNATRule struct {
	IIF []string `json:"iiface"`
	OIF []string `json:"oiface"`

	Proto []string `json:"proto"`

	SAddr    [][2]*net.IPNet `json:"saddr"`
	SAddrSet *nftables.Set   `json:"saddr_set"`

	SPort [][2]uint16 `json:"sport"`
	DPort [][2]uint16 `json:"dport"`

	DNATToIP   net.IP `json:"dnat_to_ip"`
	DNATToPort uint16 `json:"dnat_to_port"`
}
