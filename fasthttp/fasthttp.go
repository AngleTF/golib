package fasthttp

import (
	"net/http"
	"net/url"
	"fmt"
	"io/ioutil"
	"strings"
	"encoding/json"
	"time"
	"golib/fastcheck"
	"golib/fastslice"
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
		UserAgent:    "FastHttp/1.0",
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
	if !fastcheck.IsEmpty(ctx.Cookie) {
		ctx.Cookie += "; "
	}
	ctx.Cookie += cookie.String()
	return ctx
}

func (ctx *Setting) PushQueue() {
	//Dependency injection
	NewFastHttp(ctx).PushQueue()
}

func NewFastHttp(setting *Setting) *FastHttp {

	defer func() {
		if err := recover(); err != nil {
			setting.ErrorChannel <- err.(HttpError)
			return
		}
	}()

	var req *http.Request
	var body string

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
			panic(HttpError(err.Error()))
		}
		body = string(params)
	default:
		panic(HttpError("please select http request method"))
	}

	req, err = http.NewRequest(setting.Method, setting.Addr, strings.NewReader(body))
	if err != nil {
		panic(HttpError(err.Error()))
	}

	if setting.UserAgent != "" {
		req.Header.Set("User-Agent", setting.UserAgent)
	}

	if setting.Cookie != "" {
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

func serviceRequest(fastHttp *FastHttp, callback HttpRes) {
	var resp *http.Response
	var body []byte

	defer func() {
		if err := recover(); err != nil {
			fastHttp.Setting.ErrorChannel <- err.(HttpError)
			return
		}
	}()

	resp, err = callback(fastHttp.Request)

	if err != nil {
		panic(HttpError(err.Error()))
	}

	defer resp.Body.Close()

	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(HttpError(err.Error()))
	}

	err = stateAudit(resp.StatusCode, resp.Status)
	if err != nil {
		panic(err)
	}

	fastHttp.Response = resp

	fastHttp.Setting.DataChannel <- string(body)
}

func stateAudit(code int, message string) error {
	switch true {
	case code == 301 || code == 302 || code == 303 || code == 307:
		return nil
	case code >= 200 && code <= 299:
		return nil
	default:
		return HttpError(message)
	}
}

func Run() {
	client := &http.Client{}

	for lastRequestVal, flag := fastslice.Pop(&requestQueue); flag; {

		lastRequest, ok := lastRequestVal.Interface().(*FastHttp)
		if !ok{
			continue
		}

		go serviceRequest(lastRequest, func(request *http.Request) (*http.Response, error) {
			return client.Do(request)
		})

		lastRequestVal, flag = fastslice.Pop(&requestQueue)
	}
}

