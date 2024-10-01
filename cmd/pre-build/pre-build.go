package main

import (
	"encoding/hex"
	"flag"
	"fmt"
	"os"
	"sort"
	"strings"

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

// Fetches ip range sources & generates the ip range data file & tree data file
func main() {
	var writeRangesDir string
	var force bool
	flag.StringVar(&writeRangesDir, "write-ranges", "", "optionnaly store the ranges in a directory")
	flag.BoolVar(&force, "force", false, "force to re compute tree")
	flag.Parse()

	// Fetch sourceRanges then sort
	sourceRanges := source.GetAllIPRanges(source.AllSources)
	sortRanges(sourceRanges)

	// Compute the hash of the rangesStr
	hash := computeRangesHash(sourceRanges)
	log.Info("Hash of ip ranges: %s", hash)

	// Compare to previous hash
	prevHash, err := os.ReadFile(ipRangesHashPath)
	if err != nil {
		log.Fatal("Failed to read hash", err)
	}

	if hash == string(prevHash) && !force {
		log.Info("Ip ranges have not changed (same hash), skipping")
		return
	}

	// Write new hash to disk
	err = os.WriteFile(ipRangesHashPath, []byte(hash), 0644) // nolint: mnd
	if err != nil {
		log.Fatal("Failed to write hash", err)
	}

	// Build tree
	count4 := 0
	ipv4Tree := tree.NewIPv4Tree()
	count6 := 0
	ipv6Tree := tree.NewIPv6Tree()
	for _, r := range sourceRanges {
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

	if writeRangesDir != "" {
		writeRangesToDir(ipv4Tree, ipv6Tree, writeRangesDir)
		log.Info("Wrote ranges to %s", writeRangesDir)
	}
}

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
	return fmt.Sprintf("%.1f %cB", float64(b)/float64(div), "kMGTPE"[exp])
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

func computeRangesHash(sortedRanges []*source.IPRange) string {
	h := sha256.New()
	for _, r := range sortedRanges {
		h.Write([]byte(r.String()))
	}
	hash := h.Sum(nil)
	return hex.EncodeToString(hash)
}

func sortRanges(ranges []*source.IPRange) {
	sort.Slice(ranges, func(i, j int) bool {
		return ranges[i].String() > ranges[j].String()
	})
}

// Write the ranges per provider under the given directory
func writeRangesToDir(ipv4tree tree.Tree, ipv6tree tree.Tree, rangesDir string) {
	// Extract ranges from trees
	v4ranges := ipv4tree.GetAllRanges()
	v6ranges := ipv6tree.GetAllRanges()

	sortRanges(v4ranges)
	sortRanges(v6ranges)

	// Check if ranges dir exists, if not create it
	if _, err := os.Stat(rangesDir); os.IsNotExist(err) {
		err = os.Mkdir(rangesDir, os.ModePerm)
		if err != nil {
			log.Fatal("Failed to create ranges dir", err)
		}
	}

	// Map ranges per provider, first v4 then v6
	rangesPerProvider := make(map[string][]*source.IPRange)
	for _, r := range append(v4ranges, v6ranges...) {
		providerKey := strings.ToLower(r.Provider.String())
		rangesPerProvider[providerKey] = append(rangesPerProvider[providerKey], r)
	}

	// Write ranges to files
	for provider, ranges := range rangesPerProvider {
		fileContents := strings.Builder{}
		for _, r := range ranges {
			fileContents.WriteString(r.Network.String())
			fileContents.WriteString("\n")
		}

		filePath := fmt.Sprintf("%s/%s.txt", rangesDir, provider)
		err := os.WriteFile(filePath, []byte(fileContents.String()), os.ModePerm)
		if err != nil {
			log.Fatal("Failed to write ranges file", err)
		}
	}
}
