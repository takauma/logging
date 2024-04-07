package logging

import (
	"fmt"
	"log"
	"os"
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
	l.processingWriteLog(message)
}

// Info Infoレベルでログを出力します.
func (l *Logger) Info(message ...interface{}) {
	l.level = INFO
	l.processingWriteLog(message)
}

// Warn Warnレベルでログを出力します.
func (l *Logger) Warn(message ...interface{}) {
	l.level = WARN
	l.processingWriteLog(message)
}

// Error Errorレベルでログを出力します.
func (l *Logger) Error(message ...interface{}) {
	l.level = ERROR
	l.processingWriteLog(message)
}

// processingWriteLog ログ書き込み処理を行います.
func (l *Logger) processingWriteLog(message interface{}) {
	log := l.createLog(message)
	l.outConsole(log)
	l.outFile(log)
}

// getLogStr ログ出力する文字列を作成します.
func (l *Logger) createLog(message interface{}) string {
	msg := fmt.Sprintf("%v", message)
	msg = strings.TrimPrefix(msg, "[")
	msg = strings.TrimSuffix(msg, "]")

	time := time.Now().Format("2006-01-02 15:04:05.000")
	levelStr := getLevelStr(l.level)

	// 呼び出し元メソッド.
	var execFunc string
	// エラーフラグ.
	hasErr := false

	// 最初に呼び出したメソッドを取得.
	for i := 0; ; i++ {
		pc, _, _, ok := runtime.Caller(i)

		if !ok {
			hasErr = true
			break
		}

		s := runtime.FuncForPC(pc).Name()

		if !isNotLoggerFuncPath(s) {
			pc, _, _, ok := runtime.Caller(i + 1)

			if !ok {
				hasErr = true
				break
			}

			execFunc = runtime.FuncForPC(pc).Name()
			break
		}
	}

	if hasErr {
		log.Fatal("関数名取得処理で異常が発生しました。")
	}

	return fmt.Sprintf("%s %5s %s %s\n", time, levelStr, execFunc, msg)
}

// outConsole ログをコンソールに出力します.
func (l *Logger) outConsole(log string) {
	// レベルチェック.
	if !checkLevel(l.level, l.config.levelConsole) {
		return
	}
	// コンソール出力.
	os.Stdout.WriteString(log)
}

// outFile ログをファイルに出力します.
func (l *Logger) outFile(log string) {
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
