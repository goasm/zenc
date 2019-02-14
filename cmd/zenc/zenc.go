package main

import (
	"fmt"
	"os"

	"github.com/radonlab/zenc"
	flag "github.com/spf13/pflag"
)

var (
	encrypt bool
	decrypt bool
	passwd  string
	output  string
	input   string
)

func init() {
	flag.Usage = func() {
		printErr("Usage: cipher [OPTION...] FILE")
		flag.PrintDefaults()
	}
	flag.BoolVar(&encrypt, "e", false, "encrypt file")
	flag.BoolVar(&decrypt, "d", false, "decrypt file")
	flag.StringVar(&passwd, "p", "", "password to be applied")
	flag.StringVar(&output, "o", "", "file to write output\nUse - to write to standard output")
}

func printErr(message string) {
	fmt.Fprintln(os.Stderr, message)
}

func printUsageErr(message string) {
	printErr(message)
	flag.Usage()
	os.Exit(1)
}

func process() {
	var inFile, outFile *os.File
	var err error
	if input == "-" {
		// read from standard input
		inFile = os.Stdin
	} else {
		inFile, err = os.Open(input)
		if err != nil {
			printErr(fmt.Sprintf("No such file: %s", input))
		}
		defer inFile.Close()
	}
	if output == "-" {
		// write to standard output
		outFile = os.Stdout
	} else {
		outFile, err = os.OpenFile(output, os.O_WRONLY|os.O_CREATE, 0644)
		if err != nil {
			printErr(fmt.Sprintf("No such file: %s", output))
		}
		defer outFile.Close()
	}
	switch {
	case encrypt:
		zenc.EncryptFile(inFile, outFile, passwd)
	case decrypt:
		zenc.DecryptFile(inFile, outFile, passwd)
	default:
		printUsageErr("error: missing option [-e|-d]")
	}
}

func main() {
	flag.Parse()
	if flag.NArg() < 1 {
		printUsageErr("error: missing input file")
	}
	input = flag.Arg(0)
	process()
}
