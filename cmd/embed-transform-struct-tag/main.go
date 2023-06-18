package main

import (
	"flag"
	"log"

	"github.com/utgwkk/go-transform-struct-gen/internal"
)

var dstStructPath string
var srcStructPath string

func init() {
	flag.StringVar(&srcStructPath, "src-struct", "", "Source struct path (package.name) to be transformed")
	flag.StringVar(&dstStructPath, "dst-struct", "", "Destination struct path (package.name) to transform")
}

func main() {
	flag.Parse()

	if srcStructPath == "" {
		log.Fatal("-src-struct is required")
	}
	if dstStructPath == "" {
		log.Fatal("-dst-struct is required")
	}

	dst, err := internal.ResolveStruct(dstStructPath)
	if err != nil {
		log.Fatal(err)
	}
	src, err := internal.ResolveStruct(srcStructPath)
	if err != nil {
		log.Fatal(err)
	}
	if err := internal.EmbedStructTagToSrc(dst, src); err != nil {
		log.Fatal(err)
	}
}
