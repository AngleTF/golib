package fastFile

import (
	"testing"
	"fmt"
)

func TestNewFile(t *testing.T) {
	file, err := NewDir("helloWorld")
	if err != nil{
		fmt.Println(err)
		return
	}
	file.PushFileData("abs.txt", "1234")
}
