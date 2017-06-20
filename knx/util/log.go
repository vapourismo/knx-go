// Copyright 2017 Ole Krüger.
// Licensed under the MIT license which can be found in the LICENSE file.

package util

import (
	"fmt"
	"reflect"
)

// A LogTarget is used to log certain messages.
type LogTarget interface {
	Printf(format string, args ...interface{})
}

// Logger is the log target for asynchronous and non-critical errors.
var Logger LogTarget

var longestLogger = 10

// Log sends a message to the Logger.
func Log(value interface{}, format string, args ...interface{}) {
	if Logger == nil {
		return
	}

	typ := reflect.TypeOf(value).String()

	if len(typ) > longestLogger {
		longestLogger = len(typ)
	}

	Logger.Printf(
		fmt.Sprintf("%%%ds[%%p]: %%s\n", longestLogger),
		reflect.TypeOf(value).String(), value, fmt.Sprintf(format, args...),
	)
}
