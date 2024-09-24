package tree

import (
	"bytes"
	"errors"
	"net"
	"testing"

	"escape.tech/cloudfinder/internal/log"
	source "escape.tech/cloudfinder/internal/source"
)

func makeTreeHelper() *tree {
	tree := &tree{
		Root: newNode(),
		Cat:  source.CatIPv4,
	}
	networksCdirs := []string{
		"8.8.4.0/24",
		"3.5.140.0/22",
	}

	// build the tree
	for _, networkCdir := range networksCdirs {
		_, network, err := net.ParseCIDR(networkCdir)
		if err != nil {
			log.Error("Failed to parse CIDR", err)
		}
		if network == nil {
			log.Error("Failed to parse CIDR", errors.New("network is nil"))
		}

		ipRange := &source.IPRange{
			Network: network,
			Cat:     source.CatIPv4,
		}

		tree.Add(ipRange)
	}

	return tree
}

func TestFindClosestLowerDefinedChild(t *testing.T) {
	childs := make(map[uint8]*node)
	childs[0] = nil
	childs[1] = nil
	childs[2] = &node{}
	childs[3] = nil
	childs[4] = &node{}

	tests := []struct {
		input       byte
		expected    byte
		shouldError bool
	}{
		{0, 0, true},
		{1, 0, true},
		{2, 2, false},
		{3, 2, false},
		{4, 4, false},
	}

	for _, test := range tests {
		index, err := findClosestLowerDefinedChildIndex(test.input, childs)
		if err != nil && !test.shouldError {
			t.Errorf("Expected no error result for %d", test.input)
		}
		if err == nil && test.shouldError {
			t.Errorf("Expected error for %d", test.input)
		}
		if index != test.expected {
			t.Errorf("Expected index %d, got %d", test.expected, index)
		}
	}
}

func TestNetworkIp(t *testing.T) {
	_, network, err := net.ParseCIDR("0.0.0.2/24")
	if err != nil {
		t.Errorf("Failed to parse")
	}

	if network.IP.String() != net.ParseIP("0.0.0.0").String() {
		t.Errorf("Expected network IP to be 0.0.0.0, got %s", network.IP.String())
	}
}

func TestFindClosestLowerDefinedChildFor0Defined(t *testing.T) {
	childs := make(map[uint8]*node)
	childs[0] = &node{}
	childs[1] = nil
	childs[2] = nil
	childs[3] = nil
	childs[4] = &node{}

	tests := []struct {
		input       byte
		expected    byte
		shouldError bool
	}{
		{0, 0, false},
		{1, 0, false},
		{2, 0, false},
		{3, 0, false},
		{4, 4, false},
	}

	for _, test := range tests {
		index, err := findClosestLowerDefinedChildIndex(test.input, childs)
		if err != nil && !test.shouldError {
			t.Errorf("Expected no error result for %d", test.input)
		}
		if err == nil && test.shouldError {
			t.Errorf("Expected error for %d", test.input)
		}
		if index != test.expected {
			t.Errorf("Expected index %d, got %d", test.expected, index)
		}
	}
}

func TestAddToNode(t *testing.T) {
	n := newNode()
	ipRange := &source.IPRange{}
	bytes := []byte{1, 2, 3, 4}
	n.add(bytes, ipRange)

	if n.Nodes[0] != nil {
		t.Errorf("Expected node to be nil")
	}

	leaf := n.Nodes[1].Nodes[2].Nodes[3].Nodes[4]
	if leaf == nil {
		t.Errorf("Expected leaf to be non-nil")
	}
}

func TestFindInNode(t *testing.T) {
	n := newNode()
	ipRange := &source.IPRange{}
	bytes := []byte{1, 2, 3, 4}
	n.add(bytes, ipRange)

	// should find the exact leaf
	result := n.find(bytes, false)
	if result != ipRange {
		t.Errorf("Expected result to be %v, got %v", ipRange, result)
	}

	// should find the closest lower leaf
	result = n.find([]byte{1, 2, 3, 5}, false)
	if result != ipRange {
		t.Errorf("Expected result to be %v, got %v", ipRange, result)
	}

	// should not find anything
	result = n.find([]byte{1, 2, 3, 1}, false)
	if result != nil {
		t.Errorf("Expected result to be nil, got %v", result)
	}
}

func TestIPv4Tree(t *testing.T) {
	tree := makeTreeHelper()

	tests := []struct {
		ip string
		// "" for nil
		matchingRange string
	}{
		{"0.0.0.0", ""},
		{"1.2.3.4", ""},
		{"8.8.4.0", "8.8.4.0/24"},
		{"8.8.4.51", "8.8.4.0/24"},
		{"8.8.4.254", "8.8.4.0/24"},
		{"3.5.140.0", "3.5.140.0/22"},
		{"3.5.141.2", "3.5.140.0/22"},
		{"3.5.143.254", "3.5.140.0/22"},
		{"192.178.1.1", ""},
		{"255.255.255.255", ""},
	}

	for _, test := range tests {
		ip := net.ParseIP(test.ip)
		if ip == nil {
			t.Errorf("Failed to parse %s", test.ip)
		}

		result := tree.FindIPRange(ip)
		if result == nil && test.matchingRange != "" {
			t.Errorf("Expected result for %s", test.ip)
		}
		if result != nil && test.matchingRange == "" {
			t.Errorf("Expected no result for %s", test.ip)
		}
		if result != nil && test.matchingRange != "" {
			if result.Network.String() != test.matchingRange {
				t.Errorf("Expected result %s, got %s", test.matchingRange, result.Network.String())
			}
		}
	}
}

// @todo test ipv6 tree

func TestSerializeAndLoad(t *testing.T) {
	tree := makeTreeHelper()
	buffer := bytes.NewBuffer([]byte{})
	tree.SerializeTo(buffer)

	if buffer.Len() == 0 {
		t.Errorf("Expected encoded tree to be non-empty")
	}

	tree2 := NewTreeFrom(buffer, source.CatIPv4)

	testIP := net.ParseIP("3.5.141.2")
	if tree.FindIPRange(testIP) == nil {
		t.Errorf("Expected tree1 to contain range for %s", testIP)
	}
	if tree2.FindIPRange(testIP) == nil {
		t.Errorf("Expected tree2 to contain range for %s", testIP)
	}
}

// @todo test overlapping ranges
