package logging

// checkLevel ログレベルをチェックします.
func checkLevel(level, configLevel Level) bool {
	if level < configLevel {
		return false
	}
	return true
}

// getLevelStr ログレベルを文字列として取得します.
func getLevelStr(level Level) string {
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

// isNotLoggerFuncPath ロガー関数のパスではないことを確認します.
func isNotLoggerFuncPath(funcName string) bool {
	common := "github.com/takauma/logging.(*Logger)."

	switch funcName {
	case common + "Debug":
		fallthrough
	case common + "Info":
		fallthrough
	case common + "Warn":
		fallthrough
	case common + "Error":
		return false
	default:
		return true
	}
}
