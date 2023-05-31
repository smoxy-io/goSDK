package network

import "net"

const GetOutboundIPAddr = "1.1.1.1:80"

func GetOutboundIp() (net.IP, error) {
	c, err := net.Dial("udp", GetOutboundIPAddr)

	if err != nil {
		return nil, err
	}

	defer c.Close()

	return c.LocalAddr().(*net.UDPAddr).IP, nil
}
