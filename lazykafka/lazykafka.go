package lazykafka

import (
	"io/ioutil"
	"log"
)

// Logger blah
var Logger StdLogger = log.New(ioutil.Discard, "[LazyKafka] ", log.LstdFlags)

// StdLogger is used to log error messages.
type StdLogger interface {
	Print(v ...interface{})
	Printf(format string, v ...interface{})
	Println(v ...interface{})
}

// PanicHandler blah
var PanicHandler func(interface{})

// MaxRequestSize blah
var MaxRequestSize int32 = 100 * 1024 * 1024

// MaxResponseSize blah
var MaxResponseSize int32 = 100 * 1024 * 1024
