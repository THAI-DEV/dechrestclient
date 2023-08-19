package dechrestclient

import (
	"crypto/tls"
	"fmt"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
)

type restClientType struct {
	url         string
	method      string
	timeout     time.Duration
	header      map[string]string
	body        interface{}
	showConsole bool
	proxy       string
}

var client *resty.Client
var req *resty.Request

func New(url string, method string, timeout time.Duration) *restClientType {
	return &restClientType{
		url:     url,
		method:  method,
		timeout: timeout,
	}
}

func (rcv *restClientType) SetProxy(proxy string) {
	rcv.proxy = proxy
}

func (rcv *restClientType) SetHeader(head map[string]string) {
	rcv.header = head
}

func (rcv *restClientType) SetBody(body interface{}) {
	rcv.body = body
}

func (rcv *restClientType) SetShowConsole(isShow bool) {
	rcv.showConsole = isShow
}

func (rcv *restClientType) ExecClient() (*resty.Response, error) {
	client = resty.New()

	proxy := strings.TrimSpace(rcv.proxy)
	if len(proxy) > 0 {
		client.SetProxy(proxy)
	}

	client.SetTimeout(rcv.timeout)
	client.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})

	for k, v := range rcv.header {
		client.SetHeader(k, v)
		// if rcv.showConsole {
		// 	fmt.Println("header :", k, v)
		// }
	}

	req = client.R()
	// a := "{\"id\": 11,\"name\": \"test\"}"
	req.SetBody(rcv.body)

	res, err := rcv.callRest()
	return res, err
}

func MapHeader(data []string) map[string]string {
	result := map[string]string{}
	for _, v := range data {
		s := strings.Split(v, ":")
		if len(s) == 2 {
			s[0] = strings.TrimSpace(s[0])
			s[1] = strings.TrimSpace(s[1])
			result[s[0]] = s[1]
		}
	}

	return result
}

func (rcv *restClientType) callRest() (*resty.Response, error) {
	var resp *resty.Response
	var err error
	switch rcv.method {
	case "get":
		resp, err = req.Get(rcv.url)
	case "post":
		resp, err = req.Post(rcv.url)
	case "put":
		resp, err = req.Put(rcv.url)
	case "delete":
		resp, err = req.Delete(rcv.url)
	}

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	if rcv.showConsole {
		rcv.showInfo(resp, err)
	}

	return resp, nil
}

func (rcv *restClientType) showInfo(resp *resty.Response, err error) {
	fmt.Println()
	fmt.Println("Request Info  :")
	fmt.Println("  Url         :", rcv.url)
	fmt.Println("  Method      :", rcv.method)
	fmt.Println("  Head        :", rcv.header)
	fmt.Println("  Body        :", rcv.body)

	fmt.Println()
	fmt.Println("Response Info :")
	fmt.Println("  Error       :", err)
	fmt.Println("  Status Code :", resp.StatusCode())
	fmt.Println("  Result      :", resp)
	fmt.Println()
}
