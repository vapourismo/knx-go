package knx

import (
	"io/ioutil"
	l "log"
)

// Logger is the target for
var Logger = l.New(ioutil.Discard, "", l.LstdFlags)
