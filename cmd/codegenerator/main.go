package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/valencesec/codegenerator"
)

type repeatableFlag []string
// String is an implementation of the flag.Value interface
func (i *repeatableFlag) String() string {
	return fmt.Sprintf("%v", *i)
}
// Set is an implementation of the flag.Value interface
func (i *repeatableFlag) Set(value string) error {
	*i = append(*i, value)
	return nil
}

func main() {
	inFile := flag.String("in", "", "path to input file")
	outFile := flag.String("out", "", "path to output file")
	var dirsToScan repeatableFlag
	flag.Var(&dirsToScan, "dir", "directory to scan")
	chDir := flag.String("chdir", "", "change directory before scanning")
	ext := flag.String("ext", ".code_template2", "template extension when scanning dir")
	dext := flag.String("dext", ".code_template.yaml", "data inline template in file extension when scanning dir")
	flag.Parse()

	if chDir != nil && *chDir != "" {
		err := os.Chdir(*chDir)
		if err != nil {
			panic(err)
		}
	}

	if len(dirsToScan) > 0 {
		if *inFile != "" || *outFile != "" {
			panic("-dir is mutual exclusive with -in and -out")
		}
		for _, dirToScan := range dirsToScan {
			err := codegenerator.ScanFolder(dirToScan, *ext)
			if err != nil {
				panic(err)
			}
			err = codegenerator.ScanFolder(dirToScan, *dext)
			if err != nil {
				panic(err)
			}
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
