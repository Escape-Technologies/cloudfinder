package tree

import (
	"bytes"
	"errors"
	"net"
	"slices"
	"testing"

	"github.com/Escape-Technologies/cloudfinder/internal/log"
	source "github.com/Escape-Technologies/cloudfinder/internal/source"
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

func TestNetworkIp(t *testing.T) {
	_, network, err := net.ParseCIDR("0.0.0.2/24")
	if err != nil {
		t.Errorf("Failed to parse")
	}

	if network.IP.String() != net.ParseIP("0.0.0.0").String() {
		t.Errorf("Expected network IP to be 0.0.0.0, got %s", network.IP.String())
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

func TestGetAllRanges(t *testing.T) {
	tree := makeTreeHelper()
	t1ranges := tree.GetAllRanges()
	if len(t1ranges) != 2 {
		t.Errorf("Expected tree2 to have 2 ranges")
	}
}

func TestSerializeAndLoad(t *testing.T) {
	tree := makeTreeHelper()
	buffer := bytes.NewBuffer([]byte{})
	tree.SerializeTo(buffer)

	if buffer.Len() == 0 {
		t.Errorf("Expected encoded tree to be non-empty")
	}

	tree2 := NewTreeFrom(buffer, source.CatIPv4)

	t1ranges := tree.GetAllRanges()
	t2ranges := tree2.GetAllRanges()
	if len(t1ranges) != len(t2ranges) {
		t.Errorf("Expected tree2 to have %d ranges", len(t1ranges))
	}

	testIP := net.ParseIP("3.5.141.2")
	if tree.FindIPRange(testIP) == nil {
		t.Errorf("Expected tree1 to contain range for %s", testIP)
	}
	if tree2.FindIPRange(testIP) == nil {
		t.Errorf("Expected tree2 to contain range for %s", testIP)
	}
}

func cdirsToRanges(cdirs []string) []*source.IPRange {
	ranges := make([]*source.IPRange, 0)
	for _, c := range cdirs {
		net, cat := source.ParseCIDR(c)
		ranges = append(ranges, &source.IPRange{
			Network: net,
			Cat:     cat,
		})
	}
	return ranges
}

func rangesToCdirs(ranges []*source.IPRange) []string {
	cdirs := make([]string, 0)
	for _, r := range ranges {
		cdirs = append(cdirs, r.Network.String())
	}
	return cdirs
}

func TestOverlapCases(t *testing.T) {
	tests := []struct {
		name     string
		ranges   []string
		expected []string
	}{
		{
			name:     "not starting ip",
			ranges:   []string{"1.2.3.4/24"},
			expected: []string{"1.2.3.0/24"},
		},
		{
			name:     "same networks",
			ranges:   []string{"1.2.0.0/24", "1.2.0.0/24"},
			expected: []string{"1.2.0.0/24"},
		},
		{
			name:     "different networks",
			ranges:   []string{"255.255.0.0/16", "1.2.3.0/24"},
			expected: []string{"255.255.0.0/16", "1.2.3.0/24"},
		},
		{
			name:     "overlapping networks",
			ranges:   []string{"1.2.3.0/24", "1.2.0.0/16"},
			expected: []string{"1.2.0.0/16"},
		},
		{
			name:     "overlapping networks (reversed)",
			ranges:   []string{"1.2.0.0/16", "1.2.3.0/24"},
			expected: []string{"1.2.0.0/16"},
		},
	}

	for _, tt := range tests {
		ranges := cdirsToRanges(tt.ranges)
		tree := NewIPv4Tree()
		for _, r := range ranges {
			tree.Add(r)
		}

		gotRanges := tree.GetAllRanges()
		gotCdirs := rangesToCdirs(gotRanges)
		slices.Sort(gotCdirs)
		slices.Sort(tt.expected)
		if !slices.Equal(gotCdirs, tt.expected) {
			t.Errorf("[%s] Got %+v, Expected: %+v", tt.name, gotCdirs, tt.expected)
		}
	}
}
