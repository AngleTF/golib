package fastLog

import (
	"testing"
	"fmt"
)

func TestLog_Error(t *testing.T) {
	var (
		log *Log
		err error
	)

	if log, err = NewLog("log/spider_log"); err != nil{
		fmt.Println(err)
		return
	}

	log.Error("hello")
}
