package util

import (
	"fmt"
	"log"

	"github.com/pkg/errors"
)

type StackTracer interface {
	StackTrace() errors.StackTrace
}

func GetVerboseStackTrace(depth int, st StackTracer) string {
	frames := st.StackTrace()
	if depth > 0 {
		frames = frames[:depth]
	}
	return fmt.Sprintf("%+v", frames)
}

func HandleFatalError(message string, err error) {
	if err != nil {
		log.Fatalf("%s: %v", message, err)
	}
}

func HandlePanicError(message string, err error) {
	if err != nil {
		log.Panicf("%s: %v", message, err)
	}
}
