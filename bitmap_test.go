package main

import (
	"net"
	"sync"
	"testing"
)

func TestBitmap_SetAndIsSet(t *testing.T) {
	bitmap := NewBitmap(ipV4AddressesCount)

	ip := net.ParseIP("97.71.174.4").To4()
	ipUint := ipToUint32(ip)

	bitmap.Set(ipUint)
	if !bitmap.IsSet(ipUint) {
		t.Errorf("Expected bit to be set for IP %s", ip)
	}

	anotherIP := net.ParseIP("97.71.174.5").To4()
	anotherIPUint := ipToUint32(anotherIP)
	if bitmap.IsSet(anotherIPUint) {
		t.Errorf("Did not expect bit to be set for IP %s", anotherIP)
	}
}

func TestBitmap_ConcurrentAccess(t *testing.T) {
	bitmap := NewBitmap(ipV4AddressesCount)
	ips := []string{
		"97.71.174.4",
		"97.71.173.241",
		"97.71.173.235",
		"97.71.174.4",
	}

	var waitGroup sync.WaitGroup
	waitGroup.Add(len(ips))

	for _, ipStr := range ips {
		ip := net.ParseIP(ipStr).To4()
		ipUint := ipToUint32(ip)

		go func(ipUint uint32) {
			defer waitGroup.Done()

			bitmap.Set(ipUint)
		}(ipUint)
	}

	waitGroup.Wait()

	for _, ipStr := range ips {
		ip := net.ParseIP(ipStr).To4()
		ipUint := ipToUint32(ip)

		if !bitmap.IsSet(ipUint) {
			t.Errorf("Expected bit to be set for IP %s", ipStr)
		}
	}
}

func BenchmarkBitmap_Set(b *testing.B) {
	bitmap := NewBitmap(ipV4AddressesCount)
	ip := net.ParseIP("97.71.174.4").To4()
	ipUint := ipToUint32(ip)

	for n := 0; n < b.N; n++ {
		bitmap.Set(ipUint)
	}
}

func BenchmarkBitmap_IsSet(b *testing.B) {
	bitmap := NewBitmap(ipV4AddressesCount)
	ip := net.ParseIP("97.71.174.4").To4()
	ipUint := ipToUint32(ip)
	bitmap.Set(ipUint)

	for n := 0; n < b.N; n++ {
		bitmap.IsSet(ipUint)
	}
}
