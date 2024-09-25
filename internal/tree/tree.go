package tree

import (
	"encoding/gob"
	"errors"
	"fmt"
	"io"
	"net"

	"github.com/Escape-Technologies/cloudfinder/internal/log"
	source "github.com/Escape-Technologies/cloudfinder/internal/source"
)

func findClosestLowerDefinedChildIndex(i byte, childs map[uint8]*node) (uint8, error) {
	if childs == nil {
		return 0, errors.New("no lower defined child")
	}

	if len(childs) == 0 {
		return 0, errors.New("no lower defined child")
	}

	for i := i; i > 0; i-- {
		if childs[i] != nil {
			return i, nil
		}
	}
	if childs[0] != nil {
		return 0, nil
	}
	return 0, errors.New("no lower defined child")
}

// This package implements a segement tree for ip ranges, that can be saved & loaded to/from a file at build time.
type Tree interface {
	// Find the IpRange for a given ip, returns nil if none matches.
	FindIPRange(ip net.IP) *source.IPRange

	// Add a new IpRange to the tree.
	Add(ipRange *source.IPRange)

	// Serialize the tree to a buffer.
	SerializeTo(w io.Writer)
}

type tree struct {
	Root *node        `json:"r"`
	Cat  source.IPCat `json:"c"`
}

type node struct {
	// Child Nodes
	Nodes map[uint8]*node `json:"n"`

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
		Nodes: make(map[uint8]*node),
	}
}

func (t *tree) Add(ipRange *source.IPRange) {
	if ipRange.Cat != t.Cat {
		log.Fatal("Invalid IP category", fmt.Errorf("%d != %d", ipRange.Cat, t.Cat))
	}

	if matchingRange := t.FindIPRange(ipRange.Network.IP); matchingRange != nil {
		log.Debug("[Skip] - IP range %s overlaps with existing ip range %s", ipRange.Network.String(), matchingRange.Network.String())

		if matchingRange.Provider != ipRange.Provider {
			log.Info("[Skip] - IP range %s overlaps with existing ip range %s, but has different provider %s", ipRange.Network.String(), matchingRange.Network.String(), ipRange.Provider.String())
			return
		}
		ms, _ := matchingRange.Network.Mask.Size()
		is, _ := ipRange.Network.Mask.Size()
		if ms < is {
			log.Info("[Skip] - IP range %s overlaps with existing ip range %s, but has larger mask, skipping", ipRange.Network.String(), matchingRange.Network.String())
			return
		}
	}

	var ip []byte
	if ipRange.Cat == source.CatIPv4 {
		ip = ipRange.Network.IP.To4()
	} else {
		ip = ipRange.Network.IP.To16()
	}

	t.Root.add(ip, ipRange)
}

// Add a new IpRange to the tree.
func (n *node) add(ipBytes []byte, ipRange *source.IPRange) {
	if len(ipBytes) == 0 {
		n.IPRange = ipRange
		return
	}
	_byte := ipBytes[0]
	ipBytes = ipBytes[1:]
	if n.Nodes[_byte] == nil {
		n.Nodes[_byte] = newNode()
	}
	n.Nodes[_byte].add(ipBytes, ipRange)
}

// Find the IpRange for a given ip, returns nil if none matches.
// When in fallback mode, this will try to take the highest defined IP range.
func (n *node) find(ipBytes []byte, fallbackMode bool) *source.IPRange {
	if len(ipBytes) == 0 {
		return n.IPRange
	}
	_byte := ipBytes[0]
	if fallbackMode {
		// if we are in fallback mode, we want to take the highest defined IP range
		_byte = 255
	}
	ipBytes = ipBytes[1:]
	closestChildIndex, err := findClosestLowerDefinedChildIndex(_byte, n.Nodes)
	// no lower defined child, return nil
	if err != nil {
		return nil
	}

	// enable fallback mode if we have no exact match (and not already in fallback mode)
	if !fallbackMode && closestChildIndex != _byte {
		fallbackMode = true
	}

	if n.Nodes[closestChildIndex] == nil {
		log.Fatal("Should not happen", errors.New("value at closest child index is nil"))
	}
	return n.Nodes[closestChildIndex].find(ipBytes, fallbackMode)
}

func (t *tree) findClosestIPRange(ipBytes []byte) *source.IPRange {
	return t.Root.find(ipBytes, false)
}

func (t *tree) FindIPRange(ip net.IP) *source.IPRange {
	var correctCatIP net.IP
	if t.Cat == source.CatIPv4 {
		correctCatIP = ip.To4()
	} else {
		correctCatIP = ip.To16()
	}

	closestIPRange := t.findClosestIPRange(correctCatIP)
	if closestIPRange == nil {
		return nil
	}

	if closestIPRange.Network.Contains(ip) {
		return closestIPRange
	}

	return nil
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
