package logging

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"runtime"
	"strings"
	"time"
)

// Logger ロガー構造体.
type Logger struct {
	config *Config
	level  Level
}

// NewLogger ロガーを取得します.
func NewLogger(config *Config) *Logger {
	return &Logger{config: config}
}

// Debug Debugレベルでログを出力します.
func (l *Logger) Debug(message ...interface{}) {
	l.level = DEBUG
	log := l.generateLogString(message)
	l.outputLogToConsole(log)
	l.outputLogToFile(log)
}

// Info Infoレベルでログを出力します.
func (l *Logger) Info(message ...interface{}) {
	l.level = INFO
	log := l.generateLogString(message)
	l.outputLogToConsole(log)
	l.outputLogToFile(log)
}

// Warn Warnレベルでログを出力します.
func (l *Logger) Warn(message ...interface{}) {
	l.level = WARN
	log := l.generateLogString(message)
	l.outputLogToConsole(log)
	l.outputLogToFile(log)
}

// Error Errorレベルでログを出力します.
func (l *Logger) Error(message ...interface{}) {
	l.level = ERROR
	log := l.generateLogString(message)
	l.outputLogToConsole(log)
	l.outputLogToFile(log)
}

// generateLogString ログ出力する文字列を生成します.
func (l *Logger) generateLogString(message interface{}) string {
	msg := fmt.Sprintf("%v", message)
	msg = strings.TrimPrefix(msg, "[")
	msg = strings.TrimSuffix(msg, "]")

	time := time.Now().Format("2006-01-02 15:04:05.000")
	levelStr := getLevelString(l.level)

	// 自モジュール名を取得.
	selfUnitPtr, _, _, _ := runtime.Caller(0)
	selfPCName := runtime.FuncForPC(selfUnitPtr).Name()
	moduleName := regexp.MustCompile(`(\.\(.+\))*[\./]generateLogString`).ReplaceAllString(selfPCName, "")

	// エラーフラグ.
	hasErr := false

	// ログ出力を行った関数名(フルパス)と行番号を取得.
	// 少なくとも当関数とログ出力関数で2回は関数呼び出ししているので初期値は2.
	var pcName string
	var line int
	for i := 2; ; i++ {
		var up uintptr
		var ok bool
		up, _, line, ok = runtime.Caller(i)
		if !ok {
			hasErr = true
			break
		}
		pcName = runtime.FuncForPC(up).Name()
		if !strings.Contains(pcName, moduleName) {
			break
		}
	}

	if hasErr {
		log.Fatal("関数名取得処理で異常が発生しました。")
	}

	// レベル表示を5桁に揃える.
	for len(levelStr) != 5 {
		levelStr += " "
	}

	return fmt.Sprintf("%s %s %s:%d %s\n", time, levelStr, pcName, line, msg)
}

// outputLogToConsole ログをコンソールに出力します.
func (l *Logger) outputLogToConsole(log string) {
	// レベルチェック.
	if !checkLevel(l.level, l.config.levelConsole) {
		return
	}
	// コンソール出力.
	os.Stdout.WriteString(log)
}

// outputLogToFile ログをファイルに出力します.
func (l *Logger) outputLogToFile(log string) {
	// レベルチェック.
	if !checkLevel(l.level, l.config.levelFile) {
		return
	}
	// ファイル出力.
	file := l.getLogFile()
	file.WriteString(log)
	file.Close()
}

// getLogFile ログファイルを取得します.
func (l *Logger) getLogFile() *os.File {
	// ログファイル拡張子.
	fileType := ".log"

	// ファイルパスの整形.
	if l.config.filePath[len(l.config.filePath)-1:len(l.config.filePath)] != "/" {
		l.config.filePath += "/"
	}
	if l.config.fileName[0:1] == "/" {
		l.config.fileName = l.config.fileName[1:len(l.config.fileName)]
	}

	// ログファイルの存在チェック.
	if fileExist(l.config.filePath + l.config.fileName + fileType) {
		// ログファイルを開く.
		file, err := os.OpenFile(l.config.filePath+l.config.fileName+fileType, os.O_RDWR|os.O_APPEND, 0644)

		if err != nil {
			log.Fatal(err)
		}

		// ファイル情報を取得.
		fileInfo, err := file.Stat()

		if err != nil {
			log.Fatal(err)
		}

		// 最終更新日時を取得.
		modTime := fileInfo.ModTime()

		//本日日付を取得.
		nowDateTime := time.Now()
		today := time.Date(nowDateTime.Year(), nowDateTime.Month(), nowDateTime.Day(), 0, 0, 0, 0, time.Local)

		//ログファイルの最終更新日が本日より古い場合の処理.
		if modTime.Before(today) {
			file.Close()

			// 古いログファイルのファイル名を変更する.
			err := os.Rename(l.config.filePath+l.config.fileName+fileType,
				l.config.filePath+l.config.fileName+modTime.Format("_20060102")+fileType)

			if err != nil {
				log.Fatal(err)
			}

			//新たにログファイルを作成.
			return fileCreate(l.config.filePath + l.config.fileName + fileType)
		}

		return file
	}

	//存在しない場合はログファイルを新規作成.
	return fileCreate(l.config.filePath + l.config.fileName + fileType)
}
