package logging

import gobit "github.com/pot-code/gobit/pkg"

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
