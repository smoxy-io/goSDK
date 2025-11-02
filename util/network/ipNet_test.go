package network

import (
	"net"
	"testing"
)

func TestOverlap(t *testing.T) {
	tests := []struct {
		name    string
		a       string
		b       string
		want    bool
		wantErr bool
	}{
		{
			name: "overlapping networks - a contains b",
			a:    "192.168.0.0/16",
			b:    "192.168.1.0/24",
			want: true,
		},
		{
			name: "overlapping networks - b contains a",
			a:    "192.168.1.0/24",
			b:    "192.168.0.0/16",
			want: true,
		},
		{
			name: "identical networks",
			a:    "192.168.1.0/24",
			b:    "192.168.1.0/24",
			want: true,
		},
		{
			name: "non-overlapping networks",
			a:    "192.168.1.0/24",
			b:    "192.168.2.0/24",
			want: false,
		},
		{
			name: "adjacent networks",
			a:    "10.0.0.0/24",
			b:    "10.0.1.0/24",
			want: false,
		},
		{
			name: "ipv6 overlapping",
			a:    "2001:db8::/32",
			b:    "2001:db8:1::/48",
			want: true,
		},
		{
			name: "ipv6 non-overlapping",
			a:    "2001:db8::/32",
			b:    "2001:db9::/32",
			want: false,
		},
		{
			name: "ipv6 identical",
			a:    "fe80::/10",
			b:    "fe80::/10",
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, aNet, err := net.ParseCIDR(tt.a)
			if err != nil {
				t.Fatalf("failed to parse CIDR a: %v", err)
			}
			_, bNet, err := net.ParseCIDR(tt.b)
			if err != nil {
				t.Fatalf("failed to parse CIDR b: %v", err)
			}

			got := Overlap(aNet, bNet)
			if got != tt.want {
				t.Errorf("Overlap() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOverlapString(t *testing.T) {
	tests := []struct {
		name    string
		a       string
		b       string
		want    bool
		wantErr bool
	}{
		{
			name:    "overlapping networks - a contains b",
			a:       "192.168.0.0/16",
			b:       "192.168.1.0/24",
			want:    true,
			wantErr: false,
		},
		{
			name:    "overlapping networks - b contains a",
			a:       "192.168.1.0/24",
			b:       "192.168.0.0/16",
			want:    true,
			wantErr: false,
		},
		{
			name:    "non-overlapping networks",
			a:       "192.168.1.0/24",
			b:       "192.168.2.0/24",
			want:    false,
			wantErr: false,
		},
		{
			name:    "invalid CIDR a",
			a:       "invalid",
			b:       "192.168.1.0/24",
			want:    false,
			wantErr: true,
		},
		{
			name:    "invalid CIDR b",
			a:       "192.168.1.0/24",
			b:       "invalid",
			want:    false,
			wantErr: true,
		},
		{
			name:    "both invalid",
			a:       "invalid",
			b:       "also-invalid",
			want:    false,
			wantErr: true,
		},
		{
			name:    "ipv6 overlapping",
			a:       "2001:db8::/32",
			b:       "2001:db8:1::/48",
			want:    true,
			wantErr: false,
		},
		{
			name:    "ipv6 non-overlapping",
			a:       "2001:db8::/32",
			b:       "2001:db9::/32",
			want:    false,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := OverlapString(tt.a, tt.b)
			if (err != nil) != tt.wantErr {
				t.Errorf("OverlapString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("OverlapString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNextAddressBlock(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    string
		wantErr bool
	}{
		{
			name:    "ipv4 /24 network",
			input:   "192.168.1.0/24",
			want:    "192.168.2.0/24",
			wantErr: false,
		},
		{
			name:    "ipv4 /16 network",
			input:   "192.168.0.0/16",
			want:    "192.169.0.0/16",
			wantErr: false,
		},
		{
			name:    "ipv4 /8 network",
			input:   "10.0.0.0/8",
			want:    "11.0.0.0/8",
			wantErr: false,
		},
		{
			name:    "ipv4 /32 single host",
			input:   "192.168.1.1/32",
			want:    "192.168.1.2/32",
			wantErr: false,
		},
		{
			name:    "ipv4 /30 small network",
			input:   "192.168.1.0/30",
			want:    "192.168.1.4/30",
			wantErr: false,
		},
		{
			name:    "ipv6 /64 network",
			input:   "2001:db8::/64",
			want:    "2001:db8:0:1::/64",
			wantErr: false,
		},
		{
			name:    "ipv6 /48 network",
			input:   "2001:db8::/48",
			want:    "2001:db8:1::/48",
			wantErr: false,
		},
		{
			name:    "ipv6 /32 network",
			input:   "2001:db8::/32",
			want:    "2001:db9::/32",
			wantErr: false,
		},
		{
			name:    "ipv6 /128 single host",
			input:   "2001:db8::1/128",
			want:    "2001:db8::2/128",
			wantErr: false,
		},
		{
			name:    "ipv4 near overflow /24",
			input:   "255.255.255.0/24",
			want:    "",
			wantErr: true,
		},
		{
			name:    "ipv4 near overflow /16",
			input:   "255.255.0.0/16",
			want:    "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, ipNet, err := net.ParseCIDR(tt.input)
			if err != nil {
				t.Fatalf("failed to parse input CIDR: %v", err)
			}

			got, err := NextAddressBlock(ipNet)
			if (err != nil) != tt.wantErr {
				t.Errorf("NextAddressBlock() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if got.String() != tt.want {
					t.Errorf("NextAddressBlock() = %v, want %v", got.String(), tt.want)
				}
				// Verify the prefix length is preserved
				wantSize, _ := ipNet.Mask.Size()
				gotSize, _ := got.Mask.Size()
				if gotSize != wantSize {
					t.Errorf("NextAddressBlock() prefix length = %d, want %d", gotSize, wantSize)
				}
			}
		})
	}
}

func TestNextAddressBlock_Sequential(t *testing.T) {
	// Test that sequential calls produce sequential networks
	tests := []struct {
		name   string
		start  string
		count  int
		expect []string
	}{
		{
			name:  "ipv4 /24 sequence",
			start: "10.0.0.0/24",
			count: 3,
			expect: []string{
				"10.0.1.0/24",
				"10.0.2.0/24",
				"10.0.3.0/24",
			},
		},
		{
			name:  "ipv4 /30 sequence",
			start: "192.168.1.0/30",
			count: 4,
			expect: []string{
				"192.168.1.4/30",
				"192.168.1.8/30",
				"192.168.1.12/30",
				"192.168.1.16/30",
			},
		},
		{
			name:  "ipv6 /64 sequence",
			start: "2001:db8::/64",
			count: 3,
			expect: []string{
				"2001:db8:0:1::/64",
				"2001:db8:0:2::/64",
				"2001:db8:0:3::/64",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, current, err := net.ParseCIDR(tt.start)
			if err != nil {
				t.Fatalf("failed to parse start CIDR: %v", err)
			}

			for i := 0; i < tt.count; i++ {
				next, err := NextAddressBlock(current)
				if err != nil {
					t.Fatalf("NextAddressBlock() iteration %d error = %v", i, err)
				}
				if next.String() != tt.expect[i] {
					t.Errorf("NextAddressBlock() iteration %d = %v, want %v", i, next.String(), tt.expect[i])
				}
				current = next
			}
		})
	}
}

func TestNextAddressBlockString(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    string
		wantErr bool
	}{
		{
			name:    "ipv4 /24 network",
			input:   "192.168.1.0/24",
			want:    "192.168.2.0/24",
			wantErr: false,
		},
		{
			name:    "ipv4 /16 network",
			input:   "192.168.0.0/16",
			want:    "192.169.0.0/16",
			wantErr: false,
		},
		{
			name:    "ipv6 /64 network",
			input:   "2001:db8::/64",
			want:    "2001:db8:0:1::/64",
			wantErr: false,
		},
		{
			name:    "invalid CIDR",
			input:   "invalid",
			want:    "",
			wantErr: true,
		},
		{
			name:    "overflow",
			input:   "255.255.255.0/24",
			want:    "",
			wantErr: true,
		},
		{
			name:    "ipv4 single host",
			input:   "10.0.0.1/32",
			want:    "10.0.0.2/32",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NextAddressBlockString(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("NextAddressBlockString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("NextAddressBlockString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNextAddressBlockString_Consistency(t *testing.T) {
	// Test that NextAddressBlockString returns the string representation of NextAddressBlock
	tests := []string{
		"192.168.1.0/24",
		"10.0.0.0/8",
		"2001:db8::/64",
		"172.16.0.0/12",
	}

	for _, network := range tests {
		t.Run(network, func(t *testing.T) {
			_, ipNet, err := net.ParseCIDR(network)
			if err != nil {
				t.Fatalf("failed to parse CIDR: %v", err)
			}

			nextBlock, err := NextAddressBlock(ipNet)
			if err != nil {
				t.Fatalf("NextAddressBlock() error = %v", err)
			}

			nextString, err := NextAddressBlockString(network)
			if err != nil {
				t.Fatalf("NextAddressBlockString() error = %v", err)
			}

			if nextBlock.String() != nextString {
				t.Errorf("NextAddressBlockString() = %v, NextAddressBlock().String() = %v", nextString, nextBlock.String())
			}
		})
	}
}

func TestNextAddressBlock_NoOverlap(t *testing.T) {
	// Verify that consecutive address blocks don't overlap
	tests := []string{
		"192.168.0.0/24",
		"10.0.0.0/16",
		"2001:db8::/64",
		"172.16.0.0/20",
	}

	for _, network := range tests {
		t.Run(network, func(t *testing.T) {
			_, ipNet, err := net.ParseCIDR(network)
			if err != nil {
				t.Fatalf("failed to parse CIDR: %v", err)
			}

			next, err := NextAddressBlock(ipNet)
			if err != nil {
				t.Fatalf("NextAddressBlock() error = %v", err)
			}

			// Check that the networks don't overlap
			if Overlap(ipNet, next) {
				t.Errorf("NextAddressBlock() produced overlapping network: %v overlaps with %v", ipNet.String(), next.String())
			}
		})
	}
}

func BenchmarkOverlap(b *testing.B) {
	_, net1, _ := net.ParseCIDR("192.168.0.0/16")
	_, net2, _ := net.ParseCIDR("192.168.1.0/24")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Overlap(net1, net2)
	}
}

func BenchmarkOverlapString(b *testing.B) {
	net1 := "192.168.0.0/16"
	net2 := "192.168.1.0/24"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = OverlapString(net1, net2)
	}
}

func BenchmarkNextAddressBlock_IPv4(b *testing.B) {
	_, ipNet, _ := net.ParseCIDR("192.168.1.0/24")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = NextAddressBlock(ipNet)
	}
}

func BenchmarkNextAddressBlock_IPv6(b *testing.B) {
	_, ipNet, _ := net.ParseCIDR("2001:db8::/64")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = NextAddressBlock(ipNet)
	}
}

func BenchmarkNextAddressBlockString(b *testing.B) {
	network := "192.168.1.0/24"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = NextAddressBlockString(network)
	}
}
