package cmtjson

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

var (
	btSpace     byte = ' '
	btStar      byte = '*'
	btSharp     byte = '#'
	btDbQuote   byte = '"'
	btSlash     byte = '/'
	btBackslash byte = '\\'
	btLineBreak byte = '\n'

	// WriteBufSize the buffer size for write, default 1MB
	WriteBufSize = 1 << 20
)

func newStore(probableSize int) (storer, error) {
	if probableSize > WriteBufSize {
		return newFileStore()
	}
	return &bytesStore{&bytes.Buffer{}}, nil
}

// RemoveJSONCommentBytes remove comment from json data
// remove "#", "//", "/* ... */" comments
func RemoveJSONCommentBytes(data []byte) []byte {
	var (
		prev byte

		inSharpCmt, inSlashCmt, inBlockCmt, inJSONStr bool
	)
	for i := 0; i < len(data); i++ {
		wt := false
		reset := false
		switch {
		case inSharpCmt, inSlashCmt:
			if data[i] == btLineBreak {
				inSharpCmt = false
				inSlashCmt = false
			}
		case inBlockCmt:
			if prev == btStar && data[i] == btSlash {
				inBlockCmt = false
				reset = true
			}
		case inJSONStr:
			if data[i] == btDbQuote && prev != btBackslash {
				inJSONStr = false
			}
			wt = true
		default:
			if data[i] == btSharp {
				inSharpCmt = true
			} else if prev == btSlash && data[i] == btSlash {
				inSlashCmt = true
			} else if prev == btSlash && data[i] == btStar {
				inBlockCmt = true
			} else {
				wt = true
				if data[i] == btDbQuote && prev != btBackslash {
					inJSONStr = true
				}
			}
		}
		// write byte
		if i > 0 && (!wt || prev == 0) {
			data[i-1] = btSpace
		}
		// block comment need reset last byte to zero
		if reset {
			prev = 0
		} else {
			prev = data[i]
		}
	}

	// write last byte
	if inSharpCmt || inSlashCmt || inBlockCmt || prev == 0 {
		data[len(data)-1] = btSpace
	}
	return data
}

// RemoveJSONComment remove comment from r which contains json data
// remove "#", "//", "/* ... */" comments
func RemoveJSONComment(r io.Reader, probableSize int) (io.ReadCloser, error) {
	store, err := newStore(probableSize)
	if err != nil {
		return nil, err
	}

	var (
		writer = bufio.NewWriterSize(store, WriteBufSize)
		frag   = make([]byte, 4096) // 4kb read size

		n    int
		prev byte

		inSharpCmt, inSlashCmt, inBlockCmt, inJSONStr bool
	)

	for {
		n, err = r.Read(frag)
		if n <= 0 {
			break
		}

		for i := 0; i < n; i++ {
			wt := false
			reset := false
			switch {
			case inSharpCmt, inSlashCmt:
				if frag[i] == btLineBreak {
					inSharpCmt = false
					inSlashCmt = false
				}
			case inBlockCmt:
				if prev == btStar && frag[i] == btSlash {
					inBlockCmt = false
					reset = true
				}
			case inJSONStr:
				if frag[i] == btDbQuote && prev != btBackslash {
					inJSONStr = false
				}
				wt = true
			default:
				if frag[i] == btSharp {
					inSharpCmt = true
				} else if prev == btSlash && frag[i] == btSlash {
					inSlashCmt = true
				} else if prev == btSlash && frag[i] == btStar {
					inBlockCmt = true
				} else {
					wt = true
					if frag[i] == btDbQuote && prev != btBackslash {
						inJSONStr = true
					}
				}
			}
			// write byte
			if wt && prev > 0 {
				err = writer.WriteByte(prev)
				if err != nil {
					break
				}
			}
			// block comment need reset last byte to zero
			if reset {
				prev = 0
			} else {
				prev = frag[i]
			}
		}

		// read error or EOF
		if err != nil {
			break
		}
	}

	if err != nil && err != io.EOF {
		writer.Flush()
		store.Close()
		return nil, err
	}

	// write last byte
	if !inSharpCmt && !inSlashCmt && !inBlockCmt && prev > 0 {
		err = writer.WriteByte(prev)
		if err != nil {
			writer.Flush()
			store.Close()
			return nil, err
		}
	}

	err = writer.Flush()
	if err != nil {
		store.Close()
	}
	store.Ready()
	return store, nil
}

// ParseFromReader parse the json data with comment from reader r
func ParseFromReader(r io.Reader, v interface{}, probableSize int) error {
	rc, err := RemoveJSONComment(r, probableSize)
	if err != nil {
		return err
	}
	defer rc.Close()
	d, err := ioutil.ReadAll(rc)
	if err != nil {
		return err
	}
	err = json.Unmarshal(d, v)
	if err != nil {
		return err
	}
	return nil
}

// ParseFromFile parse the json data with comment from file
func ParseFromFile(filename string, v interface{}) error {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	fInfo, err := f.Stat()
	if err != nil {
		return err
	}

	if int(fInfo.Size()) <= WriteBufSize {
		buf, err := ioutil.ReadAll(f)
		if err != nil {
			return err
		}
		return ParseFromBytes(buf, v)
	}

	// always file store
	return ParseFromReader(f, v, int(fInfo.Size()))
}

// ParseFromBytes parse the json data with comment from given data
// use RemoveJSONCommentBytes to improve performance
func ParseFromBytes(data []byte, v interface{}) error {
	return json.Unmarshal(RemoveJSONCommentBytes(data), v)
	//return ParseFromReader(bytes.NewReader(data), v, len(data))
}

// only for debug
func debugf(format string, args ...interface{}) {
	fmt.Printf(format, args...)
}
