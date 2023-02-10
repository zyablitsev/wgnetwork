//go:build linux
// +build linux

package ifconfig

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"strconv"
	"strings"
)

// IPAddr returns default gw iface name, gw ip address
// and wan ip address.
func IPAddr() (string, net.IP, net.IP, error) {
	// See http://man7.org/linux/man-pages/man8/route.8.html
	const file = "/proc/net/route"
	f, err := os.Open(file)
	if err != nil {
		return "", nil, nil, fmt.Errorf("can't access %s", file)
	}
	defer f.Close()

	bytes, err := ioutil.ReadAll(f)
	if err != nil {
		return "", nil, nil, fmt.Errorf("can't read %s", file)
	}

	wanIface, gatewayIP, err := parseLinuxProcNetRoute(bytes)

	iface, err := net.InterfaceByName(wanIface)
	if err != nil {
		return "", nil, nil, errors.New("can't get interface by name")
	}

	addrs, err := iface.Addrs()
	if err != nil {
		return "", nil, nil, fmt.Errorf("can't get iface addrs: %v", err)
	}

	wanIP := net.IP{}
	found := false
	for _, addr := range addrs {
		ipnet, ok := addr.(*net.IPNet)
		if !ok || ipnet.IP.DefaultMask() == nil {
			continue
		}
		wanIP = ipnet.IP.To4()
		found = true
		break
	}

	if !found {
		return "", nil, nil, errors.New("can't found public ip")
	}

	return wanIface, gatewayIP, wanIP, nil
}

func parseLinuxProcNetRoute(f []byte) (string, net.IP, error) {
	/* /proc/net/route file:
	   Iface   Destination Gateway     Flags   RefCnt  Use Metric  Mask
	   eno1    00000000    C900A8C0    0003    0   0   100 00000000    0   00
	   eno1    0000A8C0    00000000    0001    0   0   100 00FFFFFF    0   00
	*/
	const (
		sep   = "\t" // field separator
		field = 2    // field containing hex gateway address
	)
	scanner := bufio.NewScanner(bytes.NewReader(f))
	for scanner.Scan() {
		// Skip header line
		if !scanner.Scan() {
			return "", nil, errors.New("invalid linux route file")
		}

		// get field containing gateway address
		tokens := strings.Split(scanner.Text(), sep)
		if len(tokens) <= field {
			return "", nil, errors.New("invalid linux route file")
		}
		gatewayHex := "0x" + tokens[field]
		wanIface := tokens[0]

		// cast hex address to uint32
		d, _ := strconv.ParseInt(gatewayHex, 0, 64)
		d32 := uint32(d)

		// make net.IP address from uint32
		ipd32 := make(net.IP, 4)
		binary.LittleEndian.PutUint32(ipd32, d32)

		// format net.IP to dotted ipV4 string
		return wanIface, net.IP(ipd32), nil
	}
	return "", nil, errors.New("failed to parse linux route file")
}
