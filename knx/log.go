// Copyright 2017 Ole Krüger.

package knx

import (
	"fmt"
	"io/ioutil"
	l "log"
)

// Logger is the log target for asynchronous and non-critical errors.
var Logger = l.New(ioutil.Discard, "", l.LstdFlags)

func log(value interface{}, typ string, format string, args ...interface{}) {
	Logger.Printf("%10s[%p]: %s\n", typ, value, fmt.Sprintf(format, args...))
}
