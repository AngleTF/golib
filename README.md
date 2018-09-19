### fasthttp 支持同时进行大量http请求, 每次请求都会生成一个协程进行工作, 最后将结果和错误传入指定的 channel 中, fasthttp 目前支持post, get, 以及使用post传输json数据, 下面是1.0版本的使用方式

---
 example
````golang
//声明数据存放的channel和存放错误的channel
var dataCh = make(chan string, 10)
var errorCh = make(chan error, 10)

//生成fasthttp对象
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
````
