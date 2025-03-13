package network

import (
	"github.com/smoxy-io/goSDK/util/arrays"
	"net"
	"net/netip"
)

func CreateLinkLocalAddress(mac net.HardwareAddr) string {
	prefix := make([]byte, 2)
	prefix[0] = 0xfe
	prefix[1] = 0x80

	eui64 := make([]byte, 8)
	eui64[0] = mac[0] ^ 0x02 // Flip the "universal/local" bit
	eui64[1] = mac[1]
	eui64[2] = mac[2]
	eui64[3] = 0xff
	eui64[4] = 0xfe
	eui64[5] = mac[3]
	eui64[6] = mac[4]
	eui64[7] = mac[5]

	filler := arrays.Fill(6, byte(0x00))

	addr := netip.Addr{}

	if err := addr.UnmarshalBinary(append(append(prefix, filler...), eui64...)); err != nil {
		return ""
	}

	return addr.String()
}
