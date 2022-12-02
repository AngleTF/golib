package fastHttp

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// FastHttp method list
const (
	GET  = 1
	POST = 2
	JSON = 3
)

type HttpError string

func (ctx HttpError) Error() string {
	return fmt.Sprintf("Fasthttp Error: %s", string(ctx))
}

type HttpHeader map[string]string
type HttpRes func(*http.Request) (*http.Response, error)

type Setting struct {
	Addr       string      //http url
	Method     string      //http method
	MethodType int8        //defined method
	Params     interface{} //user source body params, map _struct
	Redirect   bool        //
	UserAgent  string      //user agent
	Cookie     string
	Header     HttpHeader //user source header params, map _struct
	Url        *url.URL   //Url obj
	Timeout    time.Duration
	Transport  http.RoundTripper
}

type FastHttp struct {
	Body     string //send body params
	Request  *http.Request
	Response *http.Response
	Setting  *Setting
	Error    error
}

func Get(addr string) *FastHttp {
	return NewFastHttp("GET", addr, GET, HttpHeader{"Content-Type": "text/plain"})
}

func Post(addr string) *FastHttp {
	return NewFastHttp("POST", addr, POST, HttpHeader{"Content-Type": "application/x-www-form-urlencoded"})
}

func Json(addr string) *FastHttp {
	return NewFastHttp("POST", addr, JSON, HttpHeader{"Content-Type": "application/json"})
}

func NewFastHttp(method string, addr string, methodType int8, headers HttpHeader) *FastHttp {
	return &FastHttp{
		Setting: &Setting{
			Addr:       addr,
			Method:     method,
			MethodType: methodType,
			Redirect:   true,
			UserAgent:  "FastHttp/v1.2",
			Header:     headers,
		},
	}
}

// format into k/v
// is get method, join request url ? name = tao & age = 22
func (ctx *FastHttp) SetParams(params interface{}) *FastHttp {
	ctx.Setting.Params = params
	return ctx
}

func (ctx *FastHttp) SetTimeout(timeout time.Duration) *FastHttp {
	ctx.Setting.Timeout = timeout
	return ctx
}

func (ctx *FastHttp) SetTransport(tr *http.Transport) *FastHttp {
	ctx.Setting.Transport = tr
	return ctx
}

func (ctx *FastHttp) SetUserAgent(userAgent string) *FastHttp {
	ctx.Setting.UserAgent = userAgent
	return ctx
}

func (ctx *FastHttp) SetHeader(headers HttpHeader) *FastHttp {
	for k, v := range headers {
		ctx.Setting.Header[k] = v
	}
	return ctx
}

func (ctx *FastHttp) SetCookie(name string, value string, path string, tm time.Time) *FastHttp {
	var cookie = http.Cookie{
		Name:    name,
		Value:   value,
		Path:    path,
		Expires: tm,
	}

	if ctx.Setting.Cookie != "" {
		ctx.Setting.Cookie += "; "
	}

	ctx.Setting.Cookie += cookie.String()
	return ctx
}

func (ctx *FastHttp) SetSourceCookie(c []*http.Cookie) *FastHttp {
	for _, v := range c {
		ctx.SetCookie(v.Name, v.Value, v.Path, v.Expires)
	}
	return ctx
}

// Only HTTP types can be used
func (ctx *FastHttp) SetBasicProxyAuthorization(username, password string) *FastHttp {
	basic := "Basic " + base64.StdEncoding.EncodeToString([]byte(username+":"+password))
	ctx.Setting.Header["Proxy-Authorization"] = basic
	return ctx
}

func (ctx *FastHttp) SetAuthorization(username, password string) *FastHttp {
	basic := base64.StdEncoding.EncodeToString([]byte(username + ":" + password))
	ctx.Setting.Header["Authorization"] = basic
	return ctx
}

func (ctx *FastHttp) Build() (*FastClient, error) {

	var (
		req  *http.Request
		body string
		conf = ctx.Setting
		gps  url.Values
		ok   bool
		err  error
	)

	if conf.Params != nil {
		switch conf.MethodType {
		case GET:

			if gps, ok = conf.Params.(url.Values); !ok {
				return nil, fmt.Errorf("request parameter conversion failed")
			}

			if conf.Url, err = url.Parse(conf.Addr); err != nil {
				return nil, fmt.Errorf("request address conversion failed")
			}

			if conf.Url.RawQuery != "" {
				conf.Url.RawQuery += "&"
			}

			conf.Url.RawQuery += gps.Encode()
			conf.Addr = conf.Url.String()
		case POST:
			switch conf.Params.(type) {
			case string:
				body = conf.Params.(string)
				break
			case url.Values:
				body = conf.Params.(url.Values).Encode()
			}
		case JSON:
			var jsonParams []byte
			if jsonParams, err = json.Marshal(conf.Params); err != nil {
				return nil, fmt.Errorf("request parameter conversion failed")
			}
			body = string(jsonParams)
		}
	}

	if req, err = http.NewRequest(conf.Method, conf.Addr, strings.NewReader(body)); err != nil {
		return nil, err
	}

	for k, v := range conf.Header {
		req.Header.Set(k, v)
	}

	if conf.UserAgent != "" {
		req.Header.Set("User-Agent", conf.UserAgent)
	}

	if conf.Cookie != "" {
		req.Header.Set("Cookie", conf.Cookie)
	}

	ctx.Request = req
	ctx.Body = body

	return &FastClient{
		http: ctx,
	}, nil
}

type FastClient struct {
	http *FastHttp
}

func (f *FastClient) Fetch() ([]byte, *http.Response, error) {
	var (
		conf   = f.http.Setting
		client *http.Client
		rsp    *http.Response
		err    error
		data   []byte
	)

	client = &http.Client{
		Timeout:   conf.Timeout,
		Transport: conf.Transport,
	}

	if rsp, err = client.Do(f.http.Request); err != nil {
		return nil, nil, err
	}

	defer rsp.Body.Close()

	if rsp.StatusCode == http.StatusNotFound {
		return nil, nil, fmt.Errorf("http 404")
	}

	if data, err = ioutil.ReadAll(rsp.Body); err != nil {
		return nil, nil, err
	}

	return data, rsp, nil
}
