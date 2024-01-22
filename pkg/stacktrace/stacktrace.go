package stacktrace

import "github.com/pkg/errors"

type stackTracer interface {
	StackTrace() errors.StackTrace
}

func GetStackStrace(err error) []string {
	if err == nil {
		return []string{}
	}

	frames := err.(stackTracer).StackTrace()
	stacktrace := make([]string, 0, len(frames))
	for _, frame := range frames {
		text, _ := frame.MarshalText()
		stacktrace = append(stacktrace, string(text))
	}

	return stacktrace
}
