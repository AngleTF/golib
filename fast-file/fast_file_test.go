package fastFile

import (
	"testing"
	"fmt"
	"os"
)

func TestNewFile(t *testing.T) {
	file, err := NewDir("helloWorld", os.O_WRONLY|os.O_CREATE|os.O_APPEND, os.ModePerm)
	if err != nil{
		fmt.Println(err)
		return
	}
	file.PushFileData("abs.txt", "1234")
}
