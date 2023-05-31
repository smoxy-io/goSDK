package network

import "testing"

func TestGetOutboundIp(t *testing.T) {
	ip, err := GetOutboundIp()

	if err != nil {
		t.Errorf("GetOutboundIP() returned unexpected error: %v", err)
	}

	if ip.IsLoopback() {
		t.Errorf("GetOutboundIP() returned a looback ip: %v, wanted an outbound ip", ip)
	}

	if ip.IsUnspecified() {
		t.Errorf("GetOutboundIP() returned an unspecified ip: %v", ip)
	}
}
