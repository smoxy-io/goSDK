package network

import (
	"math/big"
	"net"
	"net/netip"
	"strconv"

	"github.com/smoxy-io/goSDK/util/errors"
)

func Overlap(a, b *net.IPNet) bool {
	return a.Contains(b.IP) || b.Contains(a.IP)
}

func OverlapString(a, b string) (bool, error) {
	_, aIpNet, aErr := net.ParseCIDR(a)

	if aErr != nil {
		return false, aErr
	}

	_, bIpNet, bErr := net.ParseCIDR(b)

	if bErr != nil {
		return false, bErr
	}

	return Overlap(aIpNet, bIpNet), nil
}

func NextAddressBlock(ipNet *net.IPNet) (*net.IPNet, error) {
	prefix, pErr := netip.ParsePrefix(ipNet.String())

	if pErr != nil {
		return nil, errors.New("failed to create Prefix from IPNet: %v", pErr)
	}

	// Ensure the prefix is masked to its network start
	prefix = prefix.Masked()

	// get the network address as a byte slice
	ipBytes := prefix.Addr().AsSlice()
	maxBits := len(ipBytes) * 8

	// Convert the IP address to a big.Int for arithmetic
	ipInt := new(big.Int).SetBytes(ipBytes)

	// Calculate the size of the network block (2^(128 - bits))
	prefixLen := prefix.Bits()

	if prefixLen > maxBits || prefixLen < 0 {
		return nil, errors.New("invalid prefix length: %d", prefixLen)
	}

	// The number of host bits is maxBits - prefixLen
	hostBits := uint(maxBits - prefixLen)

	// Calculate the increment value (1 << hostBits)
	increment := new(big.Int).Lsh(big.NewInt(1), hostBits)

	// Add the increment to the current IP
	ipInt.Add(ipInt, increment)

	// Check for overflow (if ipInt exceeds the maximum ip address)
	// Max ip value is 2^maxBits - 1, which has a byte length of 16.
	if ipInt.BitLen() > maxBits {
		// Overflowed past the max address
		return nil, errors.New("overflow. no more /%d address blocks available", prefixLen)
	}

	// Convert the big.Int back to a netip.Addr
	nextIPBytes := ipInt.Bytes()

	// big.Int.Bytes() can return a slice shorter than ip address bytes length (4 or 16), so we need to pad it if necessary
	if len(nextIPBytes) < len(ipBytes) {
		paddedBytes := make([]byte, len(ipBytes))
		copy(paddedBytes[len(ipBytes)-len(nextIPBytes):], nextIPBytes)
		nextIPBytes = paddedBytes
	}

	nextAddr, ok := netip.AddrFromSlice(nextIPBytes)

	if !ok {
		// Should not happen with valid bytes
		return nil, errors.New("failed to get ip address of next address block")
	}

	// Create the new IPNet
	_, nextIpNet, ninErr := net.ParseCIDR(nextAddr.String() + "/" + strconv.Itoa(prefixLen))

	return nextIpNet, ninErr
}

func NextAddressBlockString(network string) (string, error) {
	_, ipNet, ninErr := net.ParseCIDR(network)

	if ninErr != nil {
		return "", ninErr
	}

	nextIpNet, ninErr := NextAddressBlock(ipNet)

	if ninErr != nil {
		return "", ninErr
	}

	return nextIpNet.String(), nil
}

func IsIpAddress(ip string) bool {
	_, err := netip.ParseAddr(ip)

	return err == nil
}
