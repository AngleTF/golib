package fastCommon

import (
	"testing"
	"fmt"
)
var err error

func TestZip(t *testing.T) {
	if err = Zip("D:/data/zip/1", "D:/data/zip/2/2.zip"); err != nil {
		fmt.Println(err)
	}
}

func TestUnZip(t *testing.T) {
	if err = UnZip("D:/Redis-x64-3.2.100.zip", "/data/zip/1/2"); err != nil {
		fmt.Println(err)
	}
}

func TestDecrypt(t *testing.T) {

}

func TestEncrypt(t *testing.T) {

}

func TestJoinUrl(t *testing.T) {
	fmt.Println(JoinUrl("hello/", "/word"))
}