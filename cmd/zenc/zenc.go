package main

import (
	"fmt"
	"os"

	"github.com/goasm/zenc"
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

// UsageError means there is a problem with command usage
type UsageError struct {
	message string
}

// NewUsageError creates an UsageError with a message
func NewUsageError(message string) *UsageError {
	return &UsageError{message}
}

func (e *UsageError) Error() string {
	return e.message
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

func closeFile(fp *os.File) {
	if fp != os.Stdin && fp != os.Stdout {
		if err := fp.Close(); err != nil {
			panic(err)
		}
	}
}

func cleanup() {
	if output == "-" {
		return
	}
	if _, err := os.Stat(output); err == nil {
		// output exists
		if err := os.Remove(output); err != nil {
			panic(err)
		}
	}
}

func handleError() {
	if err := recover(); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		switch err.(type) {
		case *UsageError:
			flag.Usage()
		default:
			cleanup()
		}
		os.Exit(1)
	}
}

func process() {
	if !(encrypt || decrypt) {
		panic(NewUsageError("missing option [-e|-d]"))
	}
	var err error
	src := openInput()
	dst := openOutput()
	defer func() {
		closeFile(src)
		closeFile(dst)
	}()
	switch {
	case encrypt:
		err = zenc.EncryptFile(src, dst, passwd)
	case decrypt:
		err = zenc.DecryptFile(src, dst, passwd)
	}
	if err != nil {
		panic(err)
	}
}

func main() {
	// handle errors
	defer handleError()
	// parse flags
	flag.Parse()
	if help {
		flag.Usage()
		return
	}
	if flag.NArg() < 1 {
		panic(NewUsageError("not enough arguments"))
	}
	// run main process
	input = flag.Arg(0)
	process()
}
