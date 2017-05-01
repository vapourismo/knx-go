package knx

import (
	l "log"
)

func log(args ...interface{}) {
	l.Println(args...)
}

func logf(format string, args ...interface{}) {
	l.Printf(format + "\n", args...)
}
