package mise

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime/debug"
	"strings"
)

// FailOnError log stack and fatal with given message
func FailOnError(err error, msg string) {
	if err != nil {
		stack := debug.Stack()
		log.Fatalf("%s: %s - stack:\n%s", msg, err, stack)
	}
}

// WrapErrorMsg wrap a error with given message
func WrapErrorMsg(err error, msg string) error {
	if err != nil {
		return fmt.Errorf("%s: %s", msg, err)
	}
	return nil
}

// WrapperError is an error which wrap the prev error
type WrapperError struct {
	prev   error
	Errmsg string `json:"errmsg"`
	Errno  int    `json:"errno"`
}

func (werr *WrapperError) Error() string {
	if werr == nil {
		return ""
	}
	return werr.Errmsg
}

// GetMessages get all wrappered error errmsg
func (werr *WrapperError) GetMessages() []string {
	if werr == nil {
		return nil
	}
	var msgs = make([]string, 1)
	msgs[0] = fmt.Sprintf("wrapper: %s(%d)", werr.Errmsg, werr.Errno)

	err := werr
	for err.prev != nil {
		if tmpErr, ok := err.prev.(*WrapperError); ok {
			msgs = append(msgs, fmt.Sprintf("wrapper: %s(%d)", tmpErr.Errmsg, tmpErr.Errno))
			err = tmpErr
		} else {
			msgs = append(msgs, fmt.Sprintf("origin: %s", err.prev.Error()))
			break
		}
	}
	return msgs
}

func (werr *WrapperError) String() string {
	return strings.Join(werr.GetMessages(), " >> ")
}

// WrapError wrap a error with given errmsg
func WrapError(err error, errmsg string) error {
	if err != nil {
		return &WrapperError{prev: err, Errmsg: errmsg, Errno: 0}
	}
	return nil
}

// WrapErrorNo wrap a error with given errmsg and errno
func WrapErrorNo(err error, errmsg string, errno int) error {
	if err != nil {
		return &WrapperError{prev: err, Errmsg: errmsg, Errno: errno}
	}
	return nil
}

// GetRootPath get current cmd path
func GetRootPath() string {
	p, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		panic(err)
	}
	return p
}
