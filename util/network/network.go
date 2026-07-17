package network

import "net"

const (
	GetOutboundIPAddr   = "1.1.1.1:80"
	GetOutboundIpV6Addr = "[2606:4700:4700::1111]:80"
)

func GetOutboundIp() (net.IP, error) {
	return GetOutboundIpV4()
}

func GetOutboundIpV4() (net.IP, error) {
	c, err := net.Dial("udp", GetOutboundIPAddr)

	if err != nil {
		return nil, err
	}

	defer c.Close()

	return c.LocalAddr().(*net.UDPAddr).IP, nil
}

func GetOutboundIpV6() (net.IP, error) {
	c, err := net.Dial("udp", GetOutboundIpV6Addr)

	if err != nil {
		return nil, err
	}

	defer c.Close()

	return c.LocalAddr().(*net.UDPAddr).IP, nil
}
