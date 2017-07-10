package mise

import (
	"os"
	"path/filepath"
	"reflect"
	"syscall"
)

// SetSysProcAttrPdeathsig cross platform set syscall.SysProcAttr.Pdeathsig attr
func SetSysProcAttrPdeathsig(spa *syscall.SysProcAttr) bool {
	vspa := reflect.ValueOf(spa)
	if vspa.IsNil() {
		return false
	}
	espa := vspa.Elem()
	fv := espa.FieldByName("Pdeathsig")
	if fv.IsValid() && fv.CanSet() {
		fv.Set(reflect.ValueOf(syscall.SIGTERM))
		return true
	}
	return false
}

// GetRootPath get current cmd path
func GetRootPath() string {
	p, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		panic(err)
	}
	return p
}
