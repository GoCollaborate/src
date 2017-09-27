package logger

import (
	"fmt"
	"github.com/fatih/color"
	"log"
	"os"
	"strconv"
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

func (logger *Logger) LogHeader(content string, vars ...interface{}) {
	logger.Internal.Printf("=============="+content+"==============\n", vars...)
}

func (logger *Logger) LogProgress(content string, vars ...interface{}) {
	logger.Internal.Printf(" [PROGRESS]: "+content+"\n", vars...)
}

func (logger *Logger) LogWarning(content string, vars ...interface{}) {
	logger.Internal.Printf(" [WARN]:     "+content+"\n", vars...)
}

func (logger *Logger) LogError(content string, vars ...interface{}) {
	logger.Internal.Printf(" [ERROR]:    "+content+"\n", vars...)
}

func (logger *Logger) LogNormal(content string, vars ...interface{}) {
	logger.Internal.Printf(" [NORMAL]:   "+content+"\n", vars...)
}

func LogLogo(content ...interface{}) {
	c := color.New(color.FgBlack).Add(color.Bold)
	c.Println("============================================================")
	for _, ct := range content {
		center := fmt.Sprint(ct)
		l := len(center)
		offset := (58 - l) / 2
		c.Printf("%-"+strconv.Itoa(offset)+"s %s %"+strconv.Itoa(offset)+"s\n", "=", center, "=")
	}
	c.Println("============================================================")
}

func LogHeader(content interface{}, vars ...interface{}) {
	c := color.New(color.FgBlack).Add(color.Bold)
	c.Printf(fmt.Sprint(content)+"\n", vars...)
}

func LogProgress(content interface{}, vars ...interface{}) {
	color.Cyan(now()+" [PROGRESS]: "+fmt.Sprint(content)+"\n", vars...)
}

func LogWarning(content interface{}, vars ...interface{}) {
	color.Yellow(now()+" [WARN]:     "+fmt.Sprint(content)+"\n", vars...)
}

func LogError(content interface{}, vars ...interface{}) {
	color.Red(now()+" [ERROR]:    "+fmt.Sprint(content)+"\n", vars...)
}

func LogNormal(content interface{}, vars ...interface{}) {
	color.Blue(now()+" [NORMAL]:   "+fmt.Sprint(content)+"\n", vars...)
}

func LogListPoint(content ...interface{}) {
	fmt.Printf("%38s"+"[-] %v", "", fmt.Sprint(content)+"\n")
}

func LogNormalWithPrefix(content ...interface{}) {
	LogNormal(fmt.Sprint(content...))
}

func LogErrorWithPrefix(content ...interface{}) {
	LogNormal(fmt.Sprint(content...))
}

func now() string {
	return time.Now().Format(time.RFC3339)
}
