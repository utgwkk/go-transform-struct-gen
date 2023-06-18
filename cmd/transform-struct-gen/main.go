package main

import (
	"flag"
	"io"
	"log"
	"os"

	"github.com/utgwkk/go-transform-struct-gen/internal"
)

var dstStructPath string
var srcStructPath string
var transformerType string
var transformerName string
var outPath string

func init() {
	flag.StringVar(&srcStructPath, "src-struct", "", "Source struct path (package.name) to be transformed")
	flag.StringVar(&dstStructPath, "dst-struct", "", "Destination struct path (package.name) to transform")
	flag.StringVar(&transformerType, "type", "function", "Transformer type (function, method)")
	flag.StringVar(&transformerName, "name", "", "Transformer function (or method) name")
	flag.StringVar(&outPath, "out", "-", "Output file path (default: STDOUT)")
}

func main() {
	flag.Parse()

	if srcStructPath == "" {
		log.Fatal("-src-struct is required")
	}
	if dstStructPath == "" {
		log.Fatal("-dst-struct is required")
	}
	if transformerName == "" {
		log.Fatal("-name is required")
	}

	parsedTransformerType, err := internal.TransformerTypeFromString(transformerType)
	if err != nil {
		log.Fatal(err)
	}

	opts := &internal.GenerateOption{
		TransformerType: parsedTransformerType,
		TramsformerName: transformerName,
		DestinationPath: outPath,
	}

	dst, err := internal.ResolveStruct(dstStructPath)
	if err != nil {
		log.Fatal(err)
	}
	src, err := internal.ResolveStruct(srcStructPath)
	if err != nil {
		log.Fatal(err)
	}
	code, err := internal.Generate(dst, src, opts)
	if err != nil {
		log.Fatal(err)
	}

	var w io.WriteCloser
	if outPath == "-" {
		w = os.Stdout
	} else {
		f, err := os.Create(outPath)
		if err != nil {
			log.Fatal(err)
		}
		w = f
	}
	defer w.Close()

	if _, err := io.WriteString(w, code); err != nil {
		log.Fatal(err)
	}

	if outPath != "-" {
		log.Printf("wrote %s", outPath)
	}
}
