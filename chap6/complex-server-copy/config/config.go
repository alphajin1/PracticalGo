package config

import (
	"io"
	"log"
)

// 외부 패키지에서 접근 가능하도록 대문자 설정
type AppConfig struct {
	Logger *log.Logger
}

func InitConfig(w io.Writer) AppConfig {
	return AppConfig{Logger: log.New(w, "", log.Ldate|log.Ltime|log.Lshortfile)}
}
