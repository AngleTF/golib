package fastLog

import (
	"github.com/AngleTF/golib/fast-file"
	"os"
	"fmt"
	"github.com/AngleTF/golib/fast-common"
	"time"
)

/*

	日志目录结构
	Log
	|
	|-- spider_log
	|		2018-10-15.log


	[ 2018:10:15 17:00:00 | Error]
	....


*/

type Log struct{
	fs *fastFile.Setting
}

var template = "[ %s | %s ]\n%s\n\n"

const (
	ERROR = "ERROR"
	INFO = "INFO"
	DEBUG = "DEBUG"
)

func NewLog(logDir string) (*Log ,error){
	var (
		log = Log{}
		err error
	)

	if log.fs, err = fastFile.NewDir(logDir, os.O_WRONLY|os.O_CREATE|os.O_APPEND, os.ModePerm); err != nil{
		return nil, err
	}

	return &log, nil
}

func (l *Log) Error(format string, val ...interface{}){
	l.PushData(ERROR, format, val...)
}

func (l *Log) Info(format string, val ...interface{}){
	l.PushData(INFO, format, val...)
}

func (l *Log) PushData(logType, format string, val ...interface{}){
	var data = fmt.Sprintf(template, logType, fastCommon.Date("Y/m/d H:i:s", time.Now()), fmt.Sprintf(format, val...))
	var fileName = fastCommon.Date("Y-m-d.log", time.Now())
	fmt.Println(data)
	l.fs.PushFileData(fileName, data)
}