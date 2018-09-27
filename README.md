### fasthttp 支持同时进行大量http请求, 每次请求都会生成一个协程进行工作, 最后将结果和错误传入指定的 channel 中, fasthttp 目前支持post, get, 以及使用post传输json数据, 下面是1.0版本的使用方式

---
 example
````golang
func TestGet() {

	var dataChannel = make(chan string, 10)
	var errorChannel = make(chan error, 10)
	
	//使用get请求数据
	Get("http://127.0.0.1:8080/test.php", dataChannel, errorChannel).
		SetParams(url.Values{"name[]":[]string{"tao","11111"}}).
		PushQueue()
	
	//使用post请求数据
	Post("http://127.0.0.1:8080/test.php", dataChannel, errorChannel).
		SetParams(url.Values{"name[]":[]string{"tao","22222"}}).
		PushQueue()
		
		
	//传递json数据, post请求数据
	Json("http://127.0.0.1:8080/test.php", dataChannel, errorChannel).
		SetParams(map[string]string{"name":"tao"}).
		PushQueue()

	
	Run()

	for  {
		select {
		case data := <-dataChannel:
			fmt.Println(data)
		case err := <-errorChannel:
			fmt.Println(err)

		}
	}

}
````
