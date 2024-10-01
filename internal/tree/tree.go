package tree

import (
	"encoding/gob"
	"fmt"
	"io"
	"net"

	"github.com/Escape-Technologies/cloudfinder/internal/log"
	"github.com/Escape-Technologies/cloudfinder/internal/source"
)

// This package implements a segement tree for ip ranges, that can be saved & loaded to/from a file at build time.
type Tree interface {
	// Find the IpRange for a given ip, returns nil if none matches.
	FindIPRange(ip net.IP) *source.IPRange

	// Add a new IpRange to the tree.
	Add(ipRange *source.IPRange)

	// Serialize the tree to a buffer.
	SerializeTo(w io.Writer)

	// Get all ranges stored in tree, after duduplication ...
	GetAllRanges() []*source.IPRange
}

type tree struct {
	Root *node        `json:"r"`
	Cat  source.IPCat `json:"c"`
}

type node struct {
	// Child Nodes. Size is always 2 (0 and 1), but a [2] array won't work as it will contain nil values. This is less space efficient, but serializable
	Nodes map[byte]*node `json:"n"`

	// The ip range for this node, nil if not a leaf.
	IPRange *source.IPRange `json:"i"`
}

// NewTree creates a new tree, returns the root node.
func NewIPv4Tree() Tree {
	return &tree{
		Root: newNode(),
		Cat:  source.CatIPv4,
	}
}

func NewIPv6Tree() Tree {
	return &tree{
		Root: newNode(),
		Cat:  source.CatIPv6,
	}
}

func newNode() *node {
	return &node{
		Nodes: make(map[byte]*node),
	}
}

func (t *tree) Add(ipRange *source.IPRange) {
	if ipRange.Cat != t.Cat {
		log.Fatal("Invalid IP category", fmt.Errorf("%d != %d", ipRange.Cat, t.Cat))
	}

	var ip net.IP
	if ipRange.Cat == source.CatIPv4 {
		ip = ipRange.Network.IP.To4()
	} else {
		ip = ipRange.Network.IP.To16()
	}

	prefixLen, _ := ipRange.Network.Mask.Size()

	t.Root.add(ip, prefixLen, ipRange, 0)
}

func (n *node) add(ip net.IP, prefixLen int, ipRange *source.IPRange, bitIndex int) {
	if n.IPRange != nil {
		// A larger network already exists; skip adding
		return
	}

	if bitIndex >= prefixLen {
		// Reached the node corresponding to the prefix length
		n.IPRange = ipRange
		n.Nodes = make(map[byte]*node) // Prune any subtrees
		return
	}

	// Get the bit at the current index
	byteIndex := bitIndex / 8       // nolint:mnd
	bitOffset := 7 - (bitIndex % 8) // nolint:mnd
	bit := (ip[byteIndex] >> bitOffset) & 1

	if _, ok := n.Nodes[bit]; !ok {
		n.Nodes[bit] = newNode()
	}

	n.Nodes[bit].add(ip, prefixLen, ipRange, bitIndex+1)
}

func (n *node) find(ip net.IP, bitIndex int) *source.IPRange {
	if n.IPRange != nil {
		return n.IPRange
	}

	if bitIndex >= len(ip)*8 {
		return nil
	}

	byteIndex := bitIndex / 8       // nolint:mnd
	bitOffset := 7 - (bitIndex % 8) // nolint:mnd
	bit := (ip[byteIndex] >> bitOffset) & 1

	if _, ok := n.Nodes[bit]; !ok {
		return nil
	}

	return n.Nodes[bit].find(ip, bitIndex+1)
}

func (t *tree) FindIPRange(ip net.IP) *source.IPRange {
	var correctCatIP net.IP
	if t.Cat == source.CatIPv4 {
		correctCatIP = ip.To4()
	} else {
		correctCatIP = ip.To16()
	}

	return t.Root.find(correctCatIP, 0)
}

func (n *node) walk() []*source.IPRange {
	ranges := []*source.IPRange{}

	// Safeguard, shoud not happen
	if n == nil {
		return ranges
	}

	if n.IPRange != nil {
		ranges = append(ranges, n.IPRange)
	}

	for _, leaf := range n.Nodes {
		ranges = append(ranges, leaf.walk()...)
	}

	return ranges
}

func (t *tree) GetAllRanges() []*source.IPRange {
	return t.Root.walk()
}

func (t *tree) SerializeTo(w io.Writer) {
	err := gob.NewEncoder(w).Encode(t)
	if err != nil {
		log.Fatal("Failed to serialize tree", err)
	}
}

func NewTreeFrom(r io.Reader, cat source.IPCat) Tree {
	var t Tree
	if cat == source.CatIPv4 {
		t = NewIPv4Tree()
	} else {
		t = NewIPv6Tree()
	}

	err := gob.NewDecoder(r).Decode(t)
	if err != nil {
		log.Fatal("Failed to deserialize tree", err)
	}
	return t
}
