package static

import (
	"bytes"
	_ "embed"

	"escape.tech/cloudfinder/internal/source"
	"escape.tech/cloudfinder/internal/tree"
)

//go:embed ipv4.gob
var ipv4Content []byte

//go:embed ipv6.gob
var ipv6Content []byte

func loadTreeFromBytes(b []byte, cat source.IPCat) tree.Tree {
	reader := bytes.NewReader(b)
	// log.Info("Loading IPv%d tree.", cat)
	return tree.NewTreeFrom(reader, cat)
}

func LoadIPv4Tree() tree.Tree {
	return loadTreeFromBytes(ipv4Content, source.CatIPv4)
}

func LoadIPv6Tree() tree.Tree {
	return loadTreeFromBytes(ipv6Content, source.CatIPv6)
}
