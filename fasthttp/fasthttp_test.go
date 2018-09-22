package fasthttp

import (
	"testing"
	"fmt"
	"net/url"
)

func TestGet(t *testing.T) {
	var dataChannel = make(chan string, 10)
	var errorChannel = make(chan error, 10)
	//使用get请求数据
	Get("http://121", dataChannel, errorChannel).
		SetParams(url.Values{"name[]":[]string{"tao","li"}}).
		PushQueue()

	Run()
	for {
		select {
		case data := <-dataChannel:
			fmt.Println(data)
		case err := <-errorChannel:
			fmt.Println(err)
		default:
			fmt.Println("default")
		}
	}
}

func TestPost(t *testing.T) {
	var dataChannel = make(chan string, 10)
	var errorChannel = make(chan error, 10)

	//使用post请求数据
	Post("http://127.0.0.1:8080/test.php", dataChannel, errorChannel).
		SetParams(url.Values{"name[]":[]string{"tao"}}).
		PushQueue()

	Run()
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
		SetParams(map[string]string{"name":"tao"}).
		PushQueue()

	Run()

	for {
		select {
		case data := <-dataChannel:
			fmt.Println(data)
		case err := <-errorChannel:
			fmt.Println(err)
		}
	}
}