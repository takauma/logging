package logging

// checkLevel ログレベルをチェックします.
func checkLevel(level, configLevel Level) bool {
	return !(level < configLevel)
}

// getLevelString ログレベルを文字列として取得します.
func getLevelString(level Level) string {
	switch level {
	case DEBUG:
		return "DEBUG"
	case INFO:
		return "INFO"
	case WARN:
		return "WARN"
	case ERROR:
		return "ERROR"
	default:
		return ""
	}
}
