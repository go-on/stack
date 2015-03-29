package mw

import (
	"fmt"
	"path/filepath"
	"runtime"
)

func Caller() (info string) {
	_, file, line, _ := runtime.Caller(2)
	return fmt.Sprintf("%s:%d", filepath.FromSlash(file), line)
}
