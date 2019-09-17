package ultron

import (
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"time"

	"github.com/valyala/fasthttp"
	"go.uber.org/zap"
)

type (
	// Attacker 定义一个事务、请求，需要确保该对象是Goroutine-safe
	Attacker interface {
		Name() string
		Fire() error
	}

	// HTTPPrepareFunc 构造http.Request函数，需要调用方定义，由HTTPAttacker来发送
	HTTPPrepareFunc func() (*http.Request, error)
	// HTTPResponseCheck http.Response校验函数，可由调用方自定义，如果返回error，则视为请求失败
	HTTPResponseCheck func(resp *http.Response, body []byte) error

	// HTTPAttacker http协议的Attacker实现
	HTTPAttacker struct {
		Client     CommonHTTPClient
		Prepare    HTTPPrepareFunc
		name       string
		CheckChain []HTTPResponseCheck
	}

	// FastHTTPAttacker a http attacker base on fasthttp: https://github.com/valyala/fasthttp
	FastHTTPAttacker struct {
		Client     CommonFastHTTPClient
		Prepare    FastHTTPPrepareFunc
		name       string
		CheckChain []FastHTTPResponseCheck
	}

	// FastHTTPPrepareFunc 构造fasthttp.Request请求参数
	FastHTTPPrepareFunc func(*fasthttp.Request) error

	// FastHTTPResponseCheck check fasthttp.Response
	FastHTTPResponseCheck func(*fasthttp.Response) error

	// CommonFastHTTPClient 对FastHTTPClient的抽象interface
	CommonFastHTTPClient interface {
		Do(*fasthttp.Request, *fasthttp.Response) error
		DoDeadline(*fasthttp.Request, *fasthttp.Response, time.Time) error
		DoTimeout(*fasthttp.Request, *fasthttp.Response, time.Duration) error
	}

	// CommonHTTPClient net/http Client的抽象interface
	CommonHTTPClient interface {
		Do(r *http.Request) (*http.Response, error)
	}
)

var (
	// DefaultHTTPClient 默认http.Client
	// http://tleyden.github.io/blog/2016/11/21/tuning-the-go-http-client-library-for-load-testing/
	DefaultHTTPClient = &http.Client{
		Timeout: 90 * time.Second,
		Transport: &http.Transport{
			Proxy: nil,
			DialContext: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
				DualStack: true,
			}).DialContext,
			DisableKeepAlives:     false,
			MaxIdleConns:          2000,
			MaxIdleConnsPerHost:   1000,
			IdleConnTimeout:       30 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
		},
	}

	// DefaultFastHTTPClient define the default fasthttp client use in FastHTTPAttacker
	DefaultFastHTTPClient = &fasthttp.Client{
		Name:                "ultron",
		MaxConnsPerHost:     1000,
		MaxIdleConnDuration: 30 * time.Second,
		ReadTimeout:         30 * time.Second,
		WriteTimeout:        30 * time.Second,
	}
)

// NewHTTPAttacker 创建一个新的HTTPAttacker对象
func NewHTTPAttacker(n string, p HTTPPrepareFunc, check ...HTTPResponseCheck) *HTTPAttacker {
	return &HTTPAttacker{
		Client:     DefaultHTTPClient,
		Prepare:    p,
		name:       n,
		CheckChain: check,
	}
}

// Name 返回HTTPAttacker的名称
func (ha *HTTPAttacker) Name() string {
	return ha.name
}

// Fire HTTPAttacker发起一次请求
func (ha *HTTPAttacker) Fire() error {
	if ha.Prepare == nil {
		panic("please implement Prepare() first")
	}

	req, err := ha.Prepare()
	if err != nil {
		Logger.Error("occur error on creating new http.Request object", zap.Error(err))
		return err
	}

	resp, err := ha.Client.Do(req)
	if err != nil {
		Logger.Error("occur error on sending http request", zap.Error(err))
		return err
	}

	if ha.CheckChain == nil || len(ha.CheckChain) == 0 {
		io.Copy(ioutil.Discard, resp.Body) // no checker defined, discard body
		resp.Body.Close()
		return nil
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		Logger.Error("occur error on receiving http response", zap.Error(err))
		return err
	}
	resp.Body.Close()

	for _, check := range ha.CheckChain {
		if check == nil {
			continue
		}
		if err = check(resp, body); err != nil {
			return err
		}
	}
	return nil
}

// CheckHTTPStatusCode 检查状态码是否>=400, 如果是则视为请求失败
func CheckHTTPStatusCode(resp *http.Response, body []byte) error {
	if resp.StatusCode >= http.StatusBadRequest {
		return fmt.Errorf("bad status code: %d", resp.StatusCode)
	}
	return nil
}

// CheckFastHTTPStatusCode check if status code >= 400
func CheckFastHTTPStatusCode(resp *fasthttp.Response) error {
	if code := resp.StatusCode(); code >= fasthttp.StatusBadRequest {
		return fmt.Errorf("bad status code: %d", code)
	}
	return nil
}

// NewFastHTTPAttacker return a new instance of FastHTTPAttacker
func NewFastHTTPAttacker(n string, p FastHTTPPrepareFunc, check ...FastHTTPResponseCheck) *FastHTTPAttacker {
	a := &FastHTTPAttacker{
		Client:     DefaultFastHTTPClient,
		name:       n,
		Prepare:    p,
		CheckChain: check,
	}
	return a
}

// Name return attacker's name
func (fa *FastHTTPAttacker) Name() string {
	return fa.name
}

// Fire send request and check response
func (fa *FastHTTPAttacker) Fire() error {
	if fa.Prepare == nil {
		panic("please impl Prepare() first")
	}

	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()
	defer fa.release(req, resp)

	err := fa.Prepare(req)
	if err != nil {
		return err
	}

	if err := fa.Client.Do(req, resp); err != nil {
		return err
	}

	for _, c := range fa.CheckChain {
		if c == nil {
			continue
		}
		if err := c(resp); err != nil {
			return err
		}
	}

	return nil
}

func (fa *FastHTTPAttacker) release(req *fasthttp.Request, resp *fasthttp.Response) {
	fasthttp.ReleaseRequest(req)
	fasthttp.ReleaseResponse(resp)
}
