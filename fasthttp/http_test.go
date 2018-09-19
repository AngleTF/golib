package fasthttp

import (
	"testing"
	"fmt"
)

func TestParseHttpParams(t *testing.T) {
	var dataCh = make(chan string, 10)
	var errorCh = make(chan error, 10)
	instance.Post("http://127.0.0.1:8080/test.php").
		Param(HttpParams{"name": "陶", "age": 23}).
			Send(dataCh, errorCh)

	select {
	case data := <-dataCh:
		fmt.Println(data)
	case err := <-errorCh:
		fmt.Println(err)
	}

}
