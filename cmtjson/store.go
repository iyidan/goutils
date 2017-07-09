package cmtjson

import (
	"bytes"
	"io"
	"io/ioutil"
	"os"
)

type storer interface {
	io.ReadCloser
	io.Writer
	Ready()
}

// memory store
type bytesStore struct {
	*bytes.Buffer
}

func (jbs *bytesStore) Write(p []byte) (n int, err error) {
	return jbs.Buffer.Write(p)
}

func (jbs *bytesStore) Read(p []byte) (n int, err error) {
	return jbs.Buffer.Read(p)
}

func (jbs *bytesStore) Close() error {
	jbs.Buffer.Reset()
	jbs.Buffer = nil
	return nil
}

func (jbs *bytesStore) Ready() {
	return
}

// tmp-file store
type fileStore struct {
	f *os.File
}

func newFileStore() (*fileStore, error) {
	tmpfile, err := ioutil.TempFile("", "cmtjsonFile")
	if err != nil {
		return nil, err
	}
	return &fileStore{f: tmpfile}, nil
}

func (jfs *fileStore) Write(p []byte) (n int, err error) {
	return jfs.f.Write(p)
}

func (jfs *fileStore) Read(p []byte) (n int, err error) {
	return jfs.f.Read(p)
}

func (jfs *fileStore) Close() error {
	if jfs.f != nil {
		err := jfs.f.Close()
		if err != nil {
			return err
		}
		// remove tmp file
		os.Remove(jfs.f.Name())
		jfs.f = nil
	}
	return nil
}

func (jfs *fileStore) Ready() {
	jfs.f.Seek(0, 0)
}
