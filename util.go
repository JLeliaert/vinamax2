package vinamax2

// common utilities

import (
	"log"
	"path"
	"runtime"
)

var verbose bool = true

// Remove extension from file name.
func NoExt(file string) string {
	ext := path.Ext(file)
	return file[:len(file)-len(ext)]
}

// If err != nil, trigger log.Fatal(msg, err)
func FatalErr(err interface{}) {
	_, file, line, _ := runtime.Caller(1)
	if err != nil {
		log.Fatal(file, ":", line, err)
	}
}

func Fatal(msg ...interface{}) {
	log.Fatal(msg...)
}

// Panics if err is not nil. Signals a bug.
func PanicErr(err error) {
	if err != nil {
		log.Panic(err)
	}
}

func Log(msg ...interface{}) {
	if verbose {
		log.Println(msg...)
	}
}

// Logs the error of non-nil, plus message
func LogErr(err error, msg ...interface{}) {
	if err != nil {
		log.Println(append(msg, err)...)
	}
}

func assert(test bool) {
	if !test {
		panic("assertion failed")
	}
}
