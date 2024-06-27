package main

import (
	"os"
	"testing"
)

func TestIPDataset_CountUniqueIPs(t *testing.T) {
	file, err := os.CreateTemp("", "ipdataset_test_file")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}

	defer os.Remove(file.Name())

	ips := []string{
		"97.71.174.4\n",
		"97.71.173.241\n",
		"97.71.173.235\n",
		"97.71.174.4\n",
	}

	for _, ip := range ips {
		file.WriteString(ip)
	}

	file.Close()

	ipDataset := NewIPDataset(file.Name(), 2)
	bitmap := NewBitmap(ipV4AddressesCount)

	uniqueIPs, err := ipDataset.CountUniqueIPs(bitmap)
	if err != nil {
		t.Fatalf("failed to count unique IPs: %v", err)
	}

	expectedUniqueIPs := int64(3)
	if uniqueIPs != expectedUniqueIPs {
		t.Errorf("expected %d unique IPs, got %d", expectedUniqueIPs, uniqueIPs)
	}
}
