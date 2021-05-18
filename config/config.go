package config

import "fmt"

type LogLevel string

const (
	Debug LogLevel = "debug"
	Info  LogLevel = "info"
	Warn  LogLevel = "warn"
	Error LogLevel = "error"
)

func init() {

}

func GetVersion() string {
	return "0.0.1"
}

func GetName() string {
	return "x-ui"
}

func GetListen() string {
	return ":27827"
}

func GetCertFile() string {
	return ""
}

func GetKeyFile() string {
	return ""
}

func GetLogLevel() LogLevel {
	return Debug
}

func IsDebug() bool {
	return true
}

func GetSecret() []byte {
	return []byte("")
}

func GetDBPath() string {
	return fmt.Sprintf("/etc/%s/%s.db", GetName(), GetName())
}

func GetBasePath() string {
	return "/"
}
