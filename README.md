# ZENC

[![Build Status](https://travis-ci.org/radonlab/zenc.svg?branch=master)](https://travis-ci.org/radonlab/zenc)

ZENC is a command-line tool, a library and also a file format for data encryption.

## Installation

```
go get -u github.com/radonlab/zenc/cmd/zenc
```

## Usage

### Via command

```
zenc [OPTION...] FILE
```

| Option    | Shorthand | Description            |
| --------- | --------- | ---------------------- |
| --help    | -h        | print help message     |
| --encrypt | -e        | encrypt file           |
| --decrypt | -d        | decrypt file           |
| --passwd  | -p        | password to be applied |
| --output  | -o        | file to write output   |

### Programmatically

```
go get -u github.com/radonlab/zenc
```

```go
var err error
// encryption
err = zenc.EncryptFile(reader, writer, password)
// or decryption
err = zenc.DecryptFile(reader, writer, password)
```

## License

The MIT License.
