package util

import (
	"bytes"
	"fmt"
	"runtime"
)

func Stacktrace(skip int) []byte {
	buf := new(bytes.Buffer)
	for i := skip; ; i++ {
		pc, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}

		fmt.Fprintf(buf, "%s:%d (0x%x)\n", file, line, pc)
	}
	return buf.Bytes()
}

func Stack(skip int) string {
	pc, file, line, ok := runtime.Caller(skip)
	if !ok {
		return "runtime caller cannot recover the stack information"
	}

	return fmt.Sprintf("%s:%d (0x%x)", file, line, pc)
}
