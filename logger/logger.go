package logger

import (
	"fmt"
	"github.com/fatih/color"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"sync"
	"time"
)

var path string
var file *os.File

const NORMAL = 0
const CAPITALIZED = 1
const UPPERCASE = 2
const LOWERCASE = 3

var singleton *Logger
var once sync.Once

type Logger struct {
	Internal *log.Logger
}

func NewLogger(filePath string, prefix string, clean ...bool) (*Logger, *os.File) {
	once.Do(func() {
		var (
			err error
		)
		path = filePath
		if len(clean) > 0 && clean[0] {
			file, err = os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, 0666)
		} else {
			file, err = os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		}
		if err != nil {
			file, err = os.Create(filePath)
		}

		logger := log.New(file, prefix, log.Ldate|log.Ltime|log.LUTC)
		singleton = &Logger{logger}
	})
	return singleton, file
}

func GetLoggerInstance() *Logger {
	return singleton
}

func (logger *Logger) LogHeader(content interface{}, vars ...interface{}) {
	logger.Internal.Printf(fmt.Sprint(content)+"\n", vars...)
}

func (logger *Logger) LogProgress(content interface{}, vars ...interface{}) {
	logger.Internal.Printf(" [PROGRESS]: "+fmt.Sprint(content)+"\n", vars...)
}

func (logger *Logger) LogWarning(content interface{}, vars ...interface{}) {
	logger.Internal.Printf(" [WARN]:     "+fmt.Sprint(content)+"\n", vars...)
}

func (logger *Logger) LogError(content interface{}, vars ...interface{}) {
	logger.Internal.Printf(" [ERROR]:    "+fmt.Sprint(content)+"\n", vars...)
}

func (logger *Logger) LogNormal(content interface{}, vars ...interface{}) {
	logger.Internal.Printf(" [NORMAL]:   "+fmt.Sprint(content)+"\n", vars...)
}

func (logger *Logger) LogListPoint(content ...interface{}) {
	logger.Internal.Printf(" [LIST]:%5s"+"[-] %v", "", fmt.Sprint(content)+"\n")
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
	fmt.Printf(now()+" [LIST]:%5s"+"[-] %v", "", fmt.Sprint(content)+"\n")
}

func GetLogs() (string, error) {
	f, err := os.OpenFile(path, os.O_RDONLY, 0666)
	defer f.Close()
	b, err := ioutil.ReadAll(f)
	return string(b), err
}

func now() string {
	return time.Now().Format(time.RFC3339)
}
