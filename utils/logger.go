package utils

import "fmt"

// Logger インターフェースの定義
type Logger interface {
	Println(a ...interface{})
}

// デフォルトロガーの実装
type DefaultLogger struct{}

func (l *DefaultLogger) Println(a ...interface{}) {
	fmt.Println(a...)
}

// グローバルなロガー変数
var logger Logger = &DefaultLogger{}

// LogMessage はメッセージをログに出力する関数
func LogMessage(message string) {
	logger.Println(message)
}

// SetLogger はロガーを設定する関数（テスト用）
func SetLogger(l Logger) {
	logger = l
}
