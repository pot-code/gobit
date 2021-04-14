package util

import (
	"fmt"

	gobit "github.com/pot-code/gobit/pkg"
)

func GetVerboseStackTrace(depth int, st StackTracer) string {
	frames := st.StackTrace()
	if depth > 0 {
		frames = frames[:depth] // WARN: set 1 to skip empty line
	}
	return fmt.Sprintf("%+v", frames)
}

func SanitizeString(val string) string {
	if len([]rune(val)) > gobit.SanitizeStringLength {
		return string([]rune(val)[:gobit.SanitizeStringLength]) + gobit.SanitizedSuffix
	}
	return val
}

func SanitizeBytes(val []byte) []byte {
	if len(val) > gobit.SanitizeBytesLength {
		return []byte(string(val[:gobit.SanitizeBytesLength]) + gobit.SanitizedSuffix)
	}
	return val
}
