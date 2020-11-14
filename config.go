package logging

// Config ログ設定構造体.
type Config struct {
	filePath     string
	fileName     string
	levelFile    Level
	levelConsole Level
}

// NewConfig ログ設定を作成します.
func NewConfig(filePath, fileName string, levelFile, levelConsole Level) *Config {
	return &Config{
		filePath,
		fileName,
		levelFile,
		levelConsole,
	}
}
