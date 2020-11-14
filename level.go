package logging

// Level レベル.
type Level int

// ログレベル.
const (
	DEBUG Level = iota
	INFO
	WARN
	ERROR
	NONE
)
