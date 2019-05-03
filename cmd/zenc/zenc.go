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

func usageError(message string) {
	fmt.Fprintln(os.Stderr, "error:", message)
	flag.Usage()
	os.Exit(1)
}

func openInput() *os.File {
	if input == "-" {
		// read from standard input
		return os.Stdin
	}
	file, err := os.Open(input)
	if err != nil {
		panic(err)
	}
	return file
}

func openOutput() *os.File {
	if output == "-" {
		// write to standard output
		return os.Stdout
	}
	file, err := os.OpenFile(output, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		panic(err)
	}
	return file
}

func process() {
	switch {
	case encrypt:
		src := openInput()
		dst := openOutput()
		err := zenc.EncryptFile(src, dst, passwd)
		if err != nil {
			panic(err)
		}
		if src != os.Stdin {
			src.Close()
		}
		if dst != os.Stdout {
			dst.Close()
		}
	case decrypt:
		src := openInput()
		dst := openOutput()
		err := zenc.DecryptFile(src, dst, passwd)
		if err != nil {
			panic(err)
		}
		if src != os.Stdin {
			src.Close()
		}
		if dst != os.Stdout {
			dst.Close()
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
	// handle errors
	defer func() {
		if err := recover(); err != nil {
			fmt.Fprintln(os.Stderr, "error:", err)
			os.Exit(1)
		}
	}()
	if flag.NArg() < 1 {
		usageError("missing input file")
	}
	input = flag.Arg(0)
	process()
}
