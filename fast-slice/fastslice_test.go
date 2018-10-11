package fastSlice

import (
	"testing"
	"fmt"
)

var nameList = []int{}

func TestPop(t *testing.T) {
	lastName, eof := Pop(&nameList)
	if !eof {
		fmt.Println("弹出失败")
		return
	}
	fmt.Println(lastName.Interface().(string), nameList)
}

func TestPush(t *testing.T) {
	Push(&nameList, 111)
	fmt.Println(nameList)
}