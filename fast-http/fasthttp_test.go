package fastHttp

import (
	"fmt"
	"testing"
	//"net/url"
	"net/url"
	"time"
)

func TestGet(t *testing.T) {
	var dataChannel = make(chan string, 10)
	var errorChannel = make(chan error, 10)
	//使用get请求数据
	Get("https://www.bookbao99.net/book/201811/23/id_XNjAxNjcx.html", dataChannel, errorChannel).
		//SetParams(url.Values{"name[]":[]string{"tao","11111"}}).
		PushQueue()

	NewClient().SetTimeout(time.Second * 100).Run()

	for {
		select {
		case data := <-dataChannel:
			fmt.Println(data)
		case err := <-errorChannel:
			fmt.Println(err)
		}
	}

}

func TestPost(t *testing.T) {
	var dataChannel = make(chan string, 10)
	var errorChannel = make(chan error, 10)

	//使用post请求数据
	Post("http://127.0.0.1:8080/test.php", dataChannel, errorChannel).
		SetParams(url.Values{"name[]": []string{"tao"}}).
		PushQueue()

	NewClient().Run()
	for {
		select {
		case data := <-dataChannel:
			fmt.Println(data)
		case err := <-errorChannel:
			fmt.Println(err)
		}
	}
}

func TestJson(t *testing.T) {
	var dataChannel = make(chan string, 10)
	var errorChannel = make(chan error, 10)

	//使用post请求数据
	Json("http://127.0.0.1:8080/test.php", dataChannel, errorChannel).
		SetParams(map[string]string{"name": "tao"}).
		PushQueue()

	NewClient().SetTimeout(time.Second).Run()

	for {
		select {
		case data := <-dataChannel:
			fmt.Println(data)
		case err := <-errorChannel:
			fmt.Println(err)
		}
	}
}
