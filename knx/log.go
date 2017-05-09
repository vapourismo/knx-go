package knx

import (
	"io/ioutil"
	"fmt"
	l "log"
)

// Logger is the target for
var Logger = l.New(ioutil.Discard, "", l.LstdFlags)

//
func log(value interface{}, typ string, format string, args ...interface{}) {
	Logger.Printf("%10s[%p]: %s", typ, value, fmt.Sprintf(format, args...))
}
