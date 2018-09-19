package fasthttp

import (
	"testing"
	"fmt"
)

func TestParseHttpParams(t *testing.T) {
	var dataCh = make(chan string, 10)
	var errorCh = make(chan error, 10)

	ft := new(fastHttp)

	//使用get请求数据
	ft.Get("http://127.0.0.1:8080/test.php").
		Param(HttpParams{"name": "陶", "age": 21}).
		Send(dataCh, errorCh)

	//使用post请求数据
	ft.Post("http://127.0.0.1:8080/test.php").
		Param(HttpParams{"name": "陶", "age": 22}).
		Send(dataCh, errorCh)

	//使用post传输json数据
	ft.Json("http://127.0.0.1:8080/test.php").
		Param(HttpParams{"name": "陶", "age": 23}).
		Send(dataCh, errorCh)

	for {
		select {
		case data := <-dataCh:
			fmt.Println(data)
		case err := <-errorCh:
			fmt.Println(err)
		}
	}
}
