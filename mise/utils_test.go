package mise

import (
	"errors"
	"testing"
)

func TestWrapError(t *testing.T) {
	ori := errors.New("i am origin error")
	err := WrapErrorMsg(ori, "WrapErrorMsg")
	if err.Error() != "WrapErrorMsg: "+ori.Error() {
		t.Fatal(`err.Error() != "WrapErrorMsg: "+ori.Error()`)
	}
	t.Log(err.Error())

	err = WrapError(err, "WrapError1")
	err = WrapErrorNo(err, "WrapError2", 10001)
	errmsg := "wrapper: WrapError2(10001) >> wrapper: WrapError1(0) >> origin: WrapErrorMsg: i am origin error"
	if err.(*WrapperError).String() != errmsg {
		t.Fatal(`err.(*WrapperError).String() != errmsg`)
	}
	t.Log(err.Error())
	t.Log(err.(*WrapperError).String())
	msgs := err.(*WrapperError).GetMessages()
	t.Logf("%#v", msgs)
}
