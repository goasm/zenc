package main

import (
	"fmt"
	"os"

	"github.com/radonlab/zenc"
	flag "github.com/spf13/pflag"
)

var (
	help    bool
	encrypt bool
	decrypt bool
	passwd  string
	output  string
	input   string
)

func init() {
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage: zenc [OPTION...] FILE")
		flag.PrintDefaults()
	}
	flag.BoolVarP(&help, "help", "h", false, "print help message")
	flag.BoolVarP(&encrypt, "encrypt", "e", false, "encrypt file")
	flag.BoolVarP(&decrypt, "decrypt", "d", false, "decrypt file")
	flag.StringVarP(&passwd, "passwd", "p", "", "password to be applied")
	flag.StringVarP(&output, "output", "o", "", "file to write output\nUse - to write to standard output")
}

func error(message string) {
	fmt.Fprintln(os.Stderr, "error:", message)
	os.Exit(1)
}

func usageError(message string) {
	fmt.Fprintln(os.Stderr, "error:", message)
	flag.Usage()
	os.Exit(1)
}

func prepareInput() (*os.File, bool) {
	if input == "-" {
		// read from standard input
		return os.Stdin, false
	}
	file, err := os.Open(input)
	if err != nil {
		error(fmt.Sprintf("No such file: %s", input))
	}
	return file, true
}

func prepareOutput() (*os.File, bool) {
	if output == "-" {
		// write to standard output
		return os.Stdout, false
	}
	file, err := os.OpenFile(output, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		error(fmt.Sprintf("No such file: %s", output))
	}
	return file, true
}

func process() {
	switch {
	case encrypt:
		inFile, inClosable := prepareInput()
		outFile, outClosable := prepareOutput()
		zenc.EncryptFile(inFile, outFile, passwd)
		if inClosable {
			inFile.Close()
		}
		if outClosable {
			outFile.Close()
		}
	case decrypt:
		inFile, inClosable := prepareInput()
		outFile, outClosable := prepareOutput()
		zenc.DecryptFile(inFile, outFile, passwd)
		if inClosable {
			inFile.Close()
		}
		if outClosable {
			outFile.Close()
		}
	default:
		usageError("missing option [-e|-d]")
	}
}

func main() {
	flag.Parse()
	if help {
		flag.Usage()
		return
	}
	if flag.NArg() < 1 {
		usageError("missing input file")
	}
	input = flag.Arg(0)
	process()
}
