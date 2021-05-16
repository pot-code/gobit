package logging

const (
	SanitizeStringLength = 64
	SanitizeBytesLength  = 64
	SanitizedSuffix      = "[truncated]"
)

func SanitizeString(val string) string {
	if len([]rune(val)) > SanitizeStringLength {
		return string([]rune(val)[:SanitizeStringLength]) + SanitizedSuffix
	}
	return val
}

func SanitizeBytes(val []byte) []byte {
	if len(val) > SanitizeBytesLength {
		return []byte(string(val[:SanitizeBytesLength]) + SanitizedSuffix)
	}
	return val
}
