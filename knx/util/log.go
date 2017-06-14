// Copyright 2017 Ole Krüger.

package util

import (
	"fmt"
	"io/ioutil"
	l "log"
	"reflect"
)

// Logger is the log target for asynchronous and non-critical errors.
var Logger = l.New(ioutil.Discard, "", l.LstdFlags)

var longestLogger = 10

// Log sends a message to the Logger.
func Log(value interface{}, format string, args ...interface{}) {
	typ := reflect.TypeOf(value).String()

	if len(typ) > longestLogger {
		longestLogger = len(typ)
	}

	Logger.Printf(
		fmt.Sprintf("%%%ds[%%p]: %%s\n", longestLogger),
		reflect.TypeOf(value).String(), value, fmt.Sprintf(format, args...),
	)
}
