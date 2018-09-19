package fasthttp

import (
	"net/http"
	"net/url"
	"fmt"
	"io/ioutil"
	//"io"
	"strings"
	"encoding/json"
	//"reflect"
)

//method list
const (
	GET = 1
	POST = 2
	JSON = 3
)

type HttpError string
type HttpSuccessCallback func(string)
type HttpFailureCallback func(error)
type HttpParams map[string]interface{}
type HttpRes func(string)(*http.Response, error)

type fastHttp struct {
	Addr   string
	Method int
	Url    *url.URL
	Params HttpParams
}

func (ctx HttpError) Error() string {
	return fmt.Sprintf("Http Error: %s", string(ctx))
}

func (ctx *fastHttp) Get(addr string) (*fastHttp) {
	ctx.Addr, ctx.Method = addr, GET
	return ctx
}

func (ctx *fastHttp) Post(addr string) (*fastHttp) {
	ctx.Addr, ctx.Method = addr, POST
	return ctx
}

func (ctx *fastHttp) Json(addr string) (*fastHttp) {
	ctx.Addr, ctx.Method = addr, JSON
	return ctx
}

//format into k/v
//is get method, join request url ? name = tao & age = 22
func (ctx *fastHttp) Param(params HttpParams) *fastHttp {
	ctx.Params = params
	return ctx
}

func (ctx *fastHttp) Send(dataCh chan string, errorCh chan error){

	var err error
	ctx.Url, err = url.Parse(ctx.Addr)

	if err != nil {
		errorCh <- err
		return
	}

	switch ctx.Method {
	case GET:
		oQuery := ctx.Url.Query()
		for k, v := range ctx.Params {
			oQuery.Add(k, fmt.Sprint(v))
		}
		ctx.Url.RawQuery = oQuery.Encode()
		ctx.Addr = ctx.Url.String()
		go serviceRequest(ctx.Addr,dataCh, errorCh, func(addr string) (*http.Response, error) {
			return http.Get(addr)
		})
	case POST:
		body := ctx.ParseHttpParams(ctx.Params)
		go serviceRequest(ctx.Addr,dataCh, errorCh, func(addr string) (*http.Response, error) {
			return http.Post(addr, "application/x-www-form-urlencoded", strings.NewReader(body))
		})
	case JSON:
		encode, err := json.Marshal(ctx.Params)
		if err != nil {
			errorCh <- err
			return
		}
		go serviceRequest(ctx.Addr,dataCh, errorCh, func(addr string) (*http.Response, error) {
			return http.Post(addr, "application/json", strings.NewReader(string(encode)))
		})
	}
}

func (ctx *fastHttp) ParseHttpParams(body HttpParams) string {
	oQuery := url.Values{}
	for k, v := range body {
		oQuery.Add(k, fmt.Sprint(v))
	}
	return oQuery.Encode()
}

func serviceRequest(addr string, dataCh chan string, errCh chan error, callback HttpRes){
	resp, err := callback(addr)
	defer resp.Body.Close()
	if err != nil{
		errCh <- err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil{
		errCh <- err
	}
	dataCh <- string(body)
}

