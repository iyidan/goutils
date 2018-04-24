package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/iyidan/goutils/mise"
)

var (
	size     int
	s        string
	filename string
)

func init() {
	flag.IntVar(&size, "size", 0, "generate file size (bytes)")
	flag.StringVar(&s, "s", "", "string to repeat")
	flag.StringVar(&filename, "f", "", "output filename")
}

func main() {
	flag.Parse()
	if size <= 0 || len(s) <= 0 || len(filename) <= 0 {
		flag.Usage()
		os.Exit(1)
	}

	buf := bytes.NewBuffer(make([]byte, 0, size))
	repeatN := mise.Round(float64(size) / float64(len(s)))

	bs := []byte(s)
	for i := 0; i < repeatN; i++ {
		buf.Write(bs)
	}

	err := ioutil.WriteFile(filename, buf.Bytes(), 0600)
	if err != nil {
		fmt.Fprint(os.Stderr, err)
		os.Exit(1)
	}
}
