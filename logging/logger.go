package logger

import (
	"fmt"
	"log"
	"os"
	"strings"
)

var errorlog *os.File

const NORMAL = 0
const CAPITALIZED = 1
const UPPERCASE = 2
const LOWERCASE = 3

type Logger struct {
	Internal *log.Logger
}

func NewLogger(filePath string, prefix string) (*Logger, *os.File) {
	errorlog, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		errorlog, err = os.Create(filePath)
	}

	logger := log.New(errorlog, prefix, log.Ldate|log.Ltime|log.LUTC)
	return &Logger{logger}, errorlog
}

func (logger *Logger) LogHeader(content string, mode ...int) {
	logger.Internal.Println("=======" + transform(content, mode...) + "=======")
}

func (logger *Logger) LogProgress(content string, mode ...int) {
	logger.Internal.Println("======>" + transform(content, mode...) + "...")
}

func (logger *Logger) LogWarning(content string, mode ...int) {
	logger.Internal.Println("*******" + transform(content, mode...) + "*******")
}

func (logger *Logger) LogNormal(content string, mode ...int) {
	logger.Internal.Println(transform(content, mode...))
}

func LogHeader(content string, mode ...int) {
	fmt.Println("=======" + transform(content, mode...) + "=======")
}

func LogProgress(content string, mode ...int) {
	fmt.Println("======>" + transform(content, mode...) + "...")
}

func LogWarning(content string, mode ...int) {
	fmt.Println("*******" + transform(content, mode...) + "*******")
}

func LogNormal(content string, mode ...int) {
	fmt.Println(transform(content, mode...))
}

func transform(in string, mode ...int) (out string) {
	if len(mode) > 0 {
		switch mode[0] {
		case NORMAL:
			return in
		case CAPITALIZED:
			return strings.Title(in)
		case UPPERCASE:
			return strings.ToUpper(in)
		case LOWERCASE:
			return strings.ToLower(in)
		default:
			return in
		}
	}
	return in
}
