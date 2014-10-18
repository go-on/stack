package middleware

import (
	"fmt"
	"runtime"
)

func Caller() (info string) {
	_, file, line, _ := runtime.Caller(2)
	return fmt.Sprintf("%s:%d", file, line)
}
