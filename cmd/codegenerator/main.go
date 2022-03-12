package main

import (
	"flag"
	"github.com/tagenasec/codegenerator"
)

func main() {
	inFile := flag.String("in", "", "path to input file")
	outFile := flag.String("out", "", "path to output file")
	dirToScan := flag.String("dir", "", "directory to scan")
	ext := flag.String("ext", ".code_template2", "template extension when scanning dir")
	dext := flag.String("dext", ".code_template.yaml", "data inline template in file extension when scanning dir")
	flag.Parse()

	if *dirToScan != "" {
		if *inFile != "" || *outFile != "" {
			panic("-dir is mutual exclusive with -in and -out")
		}
		err := codegenerator.ScanFolder(*dirToScan, *ext)
		if err != nil {
			panic(err)
		}
		err = codegenerator.ScanFolder(*dirToScan, *dext)
		if err != nil {
			panic(err)
		}
	} else {
		if *inFile == "" || *outFile == "" {
			panic("either specific -dir or -in and -out")
		}
		err := codegenerator.SingleFile(*inFile, *outFile)
		if err != nil {
			panic(err)
		}
	}
}