package fastHttp

import (
	"net/http"
	"net/url"
	"fmt"
	"io/ioutil"
	"strings"
	"encoding/json"
	"time"
	"golib/fast-slice"
	"golib/fast-check"
	"crypto/tls"
	"encoding/base64"

)

//FastHttp method list
const (
	GET  = 1
	POST = 2
	JSON = 3
)

type HttpError string

func (ctx HttpError) Error() string {
	return fmt.Sprintf("Http Error: %s", string(ctx))
}

type HttpHeader map[string]string
type HttpRes func(*http.Request) (*http.Response, error)

var (
	requestQueue []*FastHttp //FastHttp queue
	err          error
)

type Setting struct {
	Addr         string      //http url
	Method       string      //http method
	MethodType   int8        //defined method
	Params       interface{} //user source body params, map struct
	Redirect     bool        //
	UserAgent    string      //user agent
	Cookie       string
	Header       HttpHeader  //user source header params, map struct
	Url          *url.URL    //Url obj
	DataChannel  chan string //response data channel
	ErrorChannel chan error  //error data channel
}

type FastHttp struct {
	Body     string //send body params
	Request  *http.Request
	Response *http.Response
	Setting  *Setting
}

func Get(addr string, dc chan string, ec chan error) (*Setting) {
	return NewSetting("GET", addr, GET, HttpHeader{"Content-Type": "text/plain"}, dc, ec)
}

func Post(addr string, dc chan string, ec chan error) (*Setting) {
	return NewSetting("POST", addr, POST, HttpHeader{"Content-Type": "application/x-www-form-urlencoded"}, dc, ec)
}

func Json(addr string, dc chan string, ec chan error) (*Setting) {
	return NewSetting("POST", addr, JSON, HttpHeader{"Content-Type": "application/json"}, dc, ec)
}

func NewSetting(method string, addr string, methodType int8, headers HttpHeader, dc chan string, ec chan error) *Setting {
	return &Setting{
		Addr:         addr,
		Method:       method,
		MethodType:   methodType,
		Redirect:     true,
		UserAgent:    "FastHttp/v1.1",
		Header:       headers,
		DataChannel:  dc,
		ErrorChannel: ec,
	}
}

//format into k/v
//is get method, join request url ? name = tao & age = 22
func (ctx *Setting) SetParams(params interface{}) *Setting {

	if gps, ok := params.(url.Values); ok {
		ctx.Url, err = url.Parse(ctx.Addr)
		if err != nil {
			ctx.ErrorChannel <- err
			return ctx
		}

		if ctx.MethodType == GET {
			ctx.Url.RawQuery = gps.Encode()
			ctx.Addr = ctx.Url.String()
		}
	}

	//fmt.Println(ctx.Addr)

	ctx.Params = params
	return ctx
}

func (ctx *Setting) SetUserAgent(userAgent string) *Setting {
	ctx.UserAgent = userAgent
	return ctx
}

func (ctx *Setting) SetHeader(headers HttpHeader) *Setting {
	for k, v := range headers {
		ctx.Header[k] = v
	}
	return ctx
}

func (ctx *Setting) SetCookie(name string, value string, path string, tm time.Time) *Setting {
	var cookie = http.Cookie{
		Name:    name,
		Value:   value,
		Path:    path,
		Expires: tm,
	}
	if !fastCheck.IsEmpty(ctx.Cookie) {
		ctx.Cookie += "; "
	}
	ctx.Cookie += cookie.String()
	return ctx
}

func (ctx *Setting) SetSourceCookie(c []*http.Cookie) *Setting {
	if c == nil{
		return ctx
	}
	for _, v := range c {
		if !fastCheck.IsEmpty(ctx.Cookie) {
			ctx.Cookie += "; "
		}
		ctx.Cookie += v.String()
	}
	return ctx
}

func (ctx *Setting) SetProxyAuthorization(username, password string) *Setting{
	ctx.Header["Proxy-Authorization"] = "Basic " + base64.StdEncoding.EncodeToString([]byte(username + ":" + password))
	return ctx
}

func (ctx *Setting) PushQueue() {
	//Dependency injection
	NewFastHttp(ctx).PushQueue()
}

func NewFastHttp(setting *Setting) *FastHttp {
	var (
		req *http.Request
		body string
	)

	switch setting.MethodType {
	case GET:
	case POST:
		if gps, ok := setting.Params.(url.Values); ok {
			body = ParseHttpParams(gps)
		} else {
			panic(HttpError("post params value type is url.Value"))
		}
	case JSON:
		var params []byte
		params, err = json.Marshal(setting.Params)
		if err != nil {
			panic(err)
		}
		body = string(params)
	default:
		panic(HttpError("please select http request method"))
	}

	if req, err = http.NewRequest(setting.Method, setting.Addr, strings.NewReader(body)); err != nil{
		panic(err)
	}

	if !fastCheck.IsEmpty(setting.UserAgent) {
		req.Header.Set("User-Agent", setting.UserAgent)
	}

	if !fastCheck.IsEmpty(setting.Cookie) {
		req.Header.Set("Cookie", setting.Cookie)
	}

	for k, v := range setting.Header {
		req.Header.Set(k, v)
	}

	return &FastHttp{
		Request: req,
		Body:    body,
		Setting: setting,
	}
}

func (ctx *FastHttp) PushQueue() {
	requestQueue = append(requestQueue, ctx)
}

func ParseHttpParams(body url.Values) string {
	return body.Encode()
}

func serviceRequest(fastHttp *FastHttp, callback HttpRes, done chan bool, queLen *int) {
	var (
		resp *http.Response
		body []byte
	)

	defer func() {
		*queLen -= 1
		if *queLen <= 0{
			done <- true
		}
	}()

	if resp, err = callback(fastHttp.Request); err != nil{
		fastHttp.Setting.ErrorChannel <- err
		return
	}

	defer func() {
		resp.Body.Close()
	}()

	if body, err = ioutil.ReadAll(resp.Body); err != nil{
		fastHttp.Setting.ErrorChannel <- err
		return
	}

	fastHttp.Response = resp

	fastHttp.Setting.DataChannel <- string(body)
}


type ClientSetting struct {
	Timeout time.Duration
	Transport *http.Transport
	Proxy func(*http.Request) (*url.URL, error)
}

func NewClient() *ClientSetting {
	return &ClientSetting{
		Timeout: time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
			DisableCompression: false,
		},
	}
}

func (ctx *ClientSetting) SetTimeout(t time.Duration) *ClientSetting {
	ctx.Timeout = t
	return ctx
}

func (ctx *ClientSetting) SetProxy(u *url.URL) *ClientSetting{
	ctx.Transport.Proxy = http.ProxyURL(u)
	return ctx
}

func (ctx *ClientSetting) SetTransport(tr *http.Transport) *ClientSetting {
	ctx.Transport = tr
	return ctx
}

func (ctx *ClientSetting) Run() chan bool {
	var (
		done = make(chan bool, 1)
		queueLen = len(requestQueue)
		lastRequest *FastHttp
		ok bool
	)

	client := &http.Client{
		Timeout: ctx.Timeout,
		Transport: ctx.Transport,
	}

	for lastRequestVal, flag := fastSlice.Pop(&requestQueue); flag; {

		if lastRequest, ok = lastRequestVal.Interface().(*FastHttp); !ok{
			continue
		}

		go serviceRequest(lastRequest, func(request *http.Request) (*http.Response, error) {
			return client.Do(request)
		}, done, &queueLen)

		lastRequestVal, flag = fastSlice.Pop(&requestQueue)
	}

	return done
}
