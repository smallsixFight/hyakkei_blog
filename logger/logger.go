package logger

import (
	"fmt"
	"github.com/smallsixFight/hyakkei_blog/util"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"
)

type Level uint8

type color uint8

const (
	DebugLevel Level = iota
	InfoLevel
	WarnLevel
	ErrorLevel
	FatalLevel
	PanicLevel
)
const (
	debugColor color = iota
	infoColor
	warnColor
	errorColor
	fatalColor
	panicColor
)

var colorNum = []color{
	debugColor: 0,
	infoColor:  32,
	warnColor:  33,
	errorColor: 31,
	fatalColor: 35,
	panicColor: 35,
}

var levelNames = []string{
	DebugLevel: "DEBUG",
	InfoLevel:  "INFO",
	WarnLevel:  "WARN",
	ErrorLevel: "ERROR",
	FatalLevel: "FATAL",
	PanicLevel: "PANIC",
}

var nextNo int
var today int64

const _baseSize = 1 << 20

type logger struct {
	mu          sync.Mutex
	isPrint     bool   // 是否在控制台输出
	openFileOut bool   // 是否启动文件保存日志信息
	outAbsPath  string // 日志文件保存路径
	outDir      string // 保存日志文件的目录
	lowestLevel Level  // 日志输出的最低级别
	fileNamePre string // 日志文件前缀
	fileMaxSize int64  // 日志文件最大大小 MB 为单位
	bufSize     int64  // 日志缓冲区大小
	showFile    bool   // 是否打印日志所在位置
}

var stdOut *logger

var lock sync.Mutex

func getLogger() *logger {
	if stdOut == nil {
		lock.Lock()
		defer lock.Unlock()
		stdOut = produce()
	}
	return stdOut
}

func (logger *logger) output(level Level, s string) {
	if logger.lowestLevel > level {
		return
	}
	now := time.Now()
	logger.mu.Lock()
	defer logger.mu.Unlock()
	logInfo := fmt.Sprintf("%s [ %s\t] %s", now.Format("2006-01-02 15:04:05"), level, s)
	var filename string
	if logger.showFile {
		_, file, line, ok := runtime.Caller(2)
		var short string
		if !ok {
			short = "???"
			line = 0
		} else if arr := strings.Split(file, string(filepath.Separator)); len(arr) > 0 {
			short = arr[len(arr)-1]
		}
		filename = fmt.Sprintf("%s:%d: ", short, line)
	}
	if logger.isPrint {
		_, _ = fmt.Fprint(os.Stdout, fmt.Sprintf("%s%c[1;0;%dm%s%c[0m", filename, 0x1B, colorNum[level], logInfo, 0x1B))
	}
	// 日志保存到文件
	if logger.openFileOut {
		file := getOutFile(logger, now)
		_, _ = file.Write([]byte(filename + logInfo))
	}
}

func getOutFile(logger *logger, now time.Time) *os.File {
	nowTime, _ := time.Parse("2006-01-02", now.Format("2006-01-02"))
	if nowTime.Unix() != today {
		nextNo = 1
		today = nowTime.Unix()
	}
	absPath := logger.outAbsPath
	if absPath == "" {
		absPath = filepath.Dir(os.Args[0])
	}
	dirPath := filepath.Join(absPath, logger.outDir)
	if !util.FileIsExist(dirPath) {
		if err := os.MkdirAll(dirPath, os.ModePerm); err != nil {
			log.Printf("create log dir fail: %s\n", err)
		}
	}
	sb := strings.Builder{}
	sb.WriteString(logger.fileNamePre)
	if sb.Len() != 0 {
		sb.WriteByte('_')
	}
	sb.WriteString(now.Format("2006-01-02"))
	baseName := sb.String()
	filename := baseName + fmt.Sprintf("_%d", nextNo) + ".log"
	filePath := filepath.Join(dirPath, filename)
	var file *os.File
	// 如果超过限制大小新建一个文件
	for {
		// 如果文件不存在
		if !util.FileIsExist(filePath) {
			file, _ := os.Create(filePath)
			return file
		}
		// 文件存在判断大小
		fileInfo, _ := os.Stat(filePath)
		if fileInfo.Size() > logger.fileMaxSize {
			nextNo++
			filename = baseName + fmt.Sprintf("_%d", nextNo) + ".log"
			filePath = filepath.Join(dirPath, filename)
		} else {
			file, _ = os.OpenFile(filePath, os.O_WRONLY|os.O_APPEND, os.ModePerm)
			return file
		}
	}
}

func Init(setters ...Option) {
	stdOut = produce(setters...)
}

func produce(setters ...Option) *logger {
	logger := &logger{
		isPrint:     true,
		openFileOut: true,
		outAbsPath:  filepath.Dir(os.Args[0]),
		outDir:      "logs",
		lowestLevel: DebugLevel,
		fileNamePre: "",
		bufSize:     2048,
		showFile:    false,
		fileMaxSize: _baseSize * 10,
	}
	for _, setter := range setters {
		setter(logger)
	}
	return logger
}

type Option func(*logger)

func LowestLevel(level Level) Option {
	return func(logger *logger) {
		logger.lowestLevel = level
	}
}

func OpenFileOut(isOpen bool) Option {
	return func(logger *logger) {
		logger.openFileOut = isOpen
	}
}

func IsPrint(v bool) Option {
	return func(logger *logger) {
		logger.isPrint = v
	}
}

func OutPath(outAbsPath string) Option {
	if !filepath.IsAbs(outAbsPath) {
		panic("outAbsPath should be an absolute path")
	}
	return func(logger *logger) {
		logger.outAbsPath = outAbsPath
	}
}

func OutDir(outDir string) Option {
	return func(logger *logger) {
		logger.outDir = outDir
	}
}

func FileNamePre(prefix string) Option {
	return func(logger *logger) {
		logger.fileNamePre = prefix
	}
}

func BufSize(size int64) Option {
	return func(logger *logger) {
		logger.bufSize = size
	}
}

func ShowFile(show bool) Option {
	return func(logger *logger) {
		logger.showFile = show
	}
}

func FileMaxSize(mb int64) Option {
	return func(logger *logger) {
		logger.fileMaxSize = _baseSize * mb
	}
}

func (l Level) String() string {
	return levelNames[l]
}

func Debug(v ...interface{}) {
	getLogger().output(DebugLevel, fmt.Sprintln(v...))
}

func Debugf(format string, v ...interface{}) {
	format += "\n"
	getLogger().output(DebugLevel, fmt.Sprintf(format, v...))
}

func Print(v ...interface{}) {
	getLogger().output(InfoLevel, fmt.Sprintln(v...))
}

func Printf(format string, v ...interface{}) {
	format += "\n"
	getLogger().output(InfoLevel, fmt.Sprintf(format, v...))
}

func Println(v ...interface{}) {
	getLogger().output(InfoLevel, fmt.Sprintln(v...))
}

func Info(v ...interface{}) {
	getLogger().output(InfoLevel, fmt.Sprintln(v...))
}

func Infof(format string, v ...interface{}) {
	format += "\n"
	getLogger().output(InfoLevel, fmt.Sprintf(format, v...))
}

func Warn(v ...interface{}) {
	getLogger().output(WarnLevel, fmt.Sprintln(v...))
}

func Warnf(format string, v ...interface{}) {
	format += "\n"
	getLogger().output(WarnLevel, fmt.Sprintf(format, v...))
}

func Error(v ...interface{}) {
	getLogger().output(ErrorLevel, fmt.Sprintln(v...))
}

func Errorf(format string, v ...interface{}) {
	format += "\n"
	getLogger().output(ErrorLevel, fmt.Sprintf(format, v...))
}

func Fatal(v ...interface{}) {
	getLogger().output(FatalLevel, fmt.Sprintln(v...))
	os.Exit(1)
}

func Fatalf(format string, v ...interface{}) {
	format += "\n"
	getLogger().output(FatalLevel, fmt.Sprintf(format, v...))
	os.Exit(1)
}

func Panic(v ...interface{}) {
	s := fmt.Sprintln(v...)
	getLogger().output(PanicLevel, s)
	panic(s)
}

func Panicf(format string, v ...interface{}) {
	format += "\n"
	s := fmt.Sprintf(format, v...)
	getLogger().output(PanicLevel, s)
	panic(s)
}
