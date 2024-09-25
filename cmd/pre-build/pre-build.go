package main

import (
	"fmt"
	"os"
	"slices"
	"sort"

	"crypto/sha256"

	"github.com/Escape-Technologies/cloudfinder/internal/log"
	"github.com/Escape-Technologies/cloudfinder/internal/source"
	"github.com/Escape-Technologies/cloudfinder/internal/tree"
)

const (
	ipv4TreePath     = "internal/static/ipv4.gob"
	ipv6TreePath     = "internal/static/ipv6.gob"
	ipRangesHashPath = "internal/static/hash.txt"
)

func byteCountSI(b int64) string {
	const unit = 1000
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB",
		float64(b)/float64(div), "kMGTPE"[exp])
}

func writeTree(tree tree.Tree, path string) {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY, 0644) // nolint: mnd
	if err != nil {
		panic(err)
	}
	// empty file
	err = f.Truncate(0)
	if err != nil {
		panic(err)
	}

	tree.SerializeTo(f)

	s, _ := f.Stat()
	size := byteCountSI(s.Size())
	log.Info("Wrote tree to %s (size: %s)", f.Name(), size)
}

// Fetches ip range sources & generates the ip range data file & tree data file
func main() {
	// Fetch ranges
	ranges := source.GetAllIPRanges(source.AllSources)

	// Sort ranges per string repr
	rangesStr := make([]string, len(ranges))
	for i, r := range ranges {
		rangesStr[i] = r.String()
	}
	sort.Strings(rangesStr)

	// Compute the hash of the rangesStr
	h := sha256.New()
	for _, r := range rangesStr {
		h.Write([]byte(r))
	}
	hash := h.Sum(nil)
	log.Info("Hash of ip ranges: %x", hash)

	// Compare to previous hash
	prevHash, err := os.ReadFile(ipRangesHashPath)
	if err != nil {
		log.Fatal("Failed to read hash", err)
	}

	if slices.Equal(hash, prevHash) {
		log.Info("Ip ranges have not changed (ip ranges hashes are the same), skipping")
		return
	}

	// Write new hash to disk
	err = os.WriteFile(ipRangesHashPath, hash, 0644) // nolint: mnd
	if err != nil {
		log.Fatal("Failed to write hash", err)
	}

	// Build tree
	count4 := 0
	ipv4Tree := tree.NewIPv4Tree()
	count6 := 0
	ipv6Tree := tree.NewIPv6Tree()
	for _, r := range ranges {
		if r.Cat == source.CatIPv4 {
			count4++
			ipv4Tree.Add(r)
		}
		if r.Cat == source.CatIPv6 {
			count6++
			ipv6Tree.Add(r)
		}
	}

	log.Info("Added %d IPv4 ranges to tree", count4)
	log.Info("Added %d IPv6 ranges to tree", count6)

	writeTree(ipv4Tree, ipv4TreePath)
	writeTree(ipv6Tree, ipv6TreePath)
}
