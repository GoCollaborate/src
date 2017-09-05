package logger

import (
	"fmt"
	"github.com/fatih/color"
	"log"
	"os"
	"strings"
	"time"
)

var errorlog *os.File

const NORMAL = 0
const CAPITALIZED = 1
const UPPERCASE = 2
const LOWERCASE = 3

type Logger struct {
	Internal *log.Logger
}

func NewLogger(filePath string, prefix string, clean ...bool) (*Logger, *os.File) {
	var (
		errorlog *os.File
		err      error
	)
	if len(clean) > 0 && clean[0] {
		errorlog, err = os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, 0666)
	} else {
		errorlog, err = os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	}
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

func (logger *Logger) LogError(content string, mode ...int) {
	logger.Internal.Println("#######" + transform(content, mode...) + "#######")
}

func (logger *Logger) LogNormal(content string, mode ...int) {
	logger.Internal.Println(transform(content, mode...))
}

func LogHeader(content interface{}, mode ...int) {
	c := color.New(color.FgBlack).Add(color.Bold)
	c.Println("=======" + transform(fmt.Sprint(content), mode...) + "=======")
}

func LogProgress(content interface{}, mode ...int) {
	color.Cyan("======>" + transform(fmt.Sprint(content), mode...) + "...")
}

func LogWarning(content interface{}, mode ...int) {
	color.Yellow(now() + ": " + transform(fmt.Sprint(content), mode...))
}

func LogError(content interface{}, mode ...int) {
	color.Red(now() + ": " + transform(fmt.Sprint(content), mode...))
}

func LogNormal(content interface{}, mode ...int) {
	fmt.Println(now() + ": " + transform(fmt.Sprint(content), mode...))
}

func LogListPoint(content interface{}, mode ...int) {
	fmt.Printf("                           " + "- " + transform(fmt.Sprint(content), mode...) + "\n")
}

func LogNormalWithPrefix(mode int, content ...interface{}) {
	LogNormal(fmt.Sprint(content...), mode)
}

func LogErrorWithPrefix(mode int, content ...interface{}) {
	LogNormal(fmt.Sprint(content...), mode)
}

func transform(in string, mode ...int) string {
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

func now() string {
	return time.Now().Format(time.RFC3339)
}
