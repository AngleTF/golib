package fastslice

import (
	"testing"
	"fmt"
)

var nameList = []interface{}{1,1,1,1}

func TestPop(t *testing.T) {
	lastName, eof := Pop(&nameList)
	if !eof {
		fmt.Println("弹出失败")
		return
	}
	fmt.Println(lastName.Interface().(string), nameList)
}

func TestPush(t *testing.T) {
	flag := Push(&nameList, "11")
	fmt.Println(nameList, flag)
}