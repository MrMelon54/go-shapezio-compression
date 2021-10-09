package main

import (
	"fmt"
	"os"

	LZString "github.com/Lazarus/lz-string-go"
	shapezio "github.com/MrMelon54/go-shapezio-compression"
)

func main() {
	fmt.Printf("Shapez.io Compression -- by MrMelon\n")

	// Download options
	if len(os.Args) == 4 {
		// List versions
		switch os.Args[1] {
		case "--compress":
		case "-c":
			shapezio.CompressFile(os.Args[2], os.Args[3])
		case "--decompress":
		case "-d":
			shapezio.DecompressFile(os.Args[2], os.Args[3])
		case "-az":
			az(os.Args[2], os.Args[3])
		}
	}

	// Help options
	if len(os.Args) == 1 {
		fmt.Printf("shapezio-compression -c <input file> <output file> - Compress the input file\n")
		fmt.Printf("shapezio-compression -d <input file> <output file> - Decompress the input file\n")
		return
	}
}

const _keyStrUriSafe = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+-$"

func az(file1 string, file2 string) {
	dat, err := os.ReadFile(file1)
	if err != nil {
		panic(err)
	}
	a, err := LZString.Decompress(string(dat[1:]), _keyStrUriSafe)
	if err != nil {
		panic(err)
	}
	err = os.WriteFile(file2, []byte(a), 0666)
	if err != nil {
		panic(err)
	}
}
