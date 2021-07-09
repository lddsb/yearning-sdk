package yearning

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

const Unauthorized = Error("Unauthorized")

type Error string

func (e Error) Error() string { return string(e) }

func (Error) YearningError() {}

type Yearning struct {
	Token    string
	username string
	password string
	host     string
}

// NewClient new a client
func NewClient(username, password, host, token string) *Yearning {
	return &Yearning{
		username: username,
		password: password,
		host:     host,
		Token: token,
	}
}

// request base request method
func (y *Yearning) request(method, path string, headers map[string]string, params map[string]interface{}) (b []byte, err error) {
	jsonData, err := json.Marshal(params)
	if err != nil {
		return
	}
	client := &http.Client{}
	req, err := http.NewRequest(method, y.host+path, bytes.NewReader(jsonData))
	if err != nil {
		return
	}
	// 处理GET请求的参数
	if strings.ToUpper(method) == "GET" {
		q := req.URL.Query()
		for key, val := range params {
			q.Add(key, val.(string))
		}
		req.URL.RawQuery = q.Encode()
	}

	// set headers
	if headers == nil {
		headers = make(map[string]string)
	}
	// 头自动带上token
	headers["Authorization"] = "Bearer " + y.Token
	headers["Content-Type"] = "application/json;charset=UTF-8"
	for headerField, headerVal := range headers {
		if req.Header.Get(headerField) == "" {
			req.Header.Add(headerField, headerVal)
		} else {
			req.Header.Set(headerField, headerVal)
		}
	}

	// do request
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %v", err)
	}

	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	b, err = io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response body failed: %v", err)
	}

	if resp.StatusCode == 401 || resp.StatusCode == 400 {
		return nil, Unauthorized
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("response error, status code: %d, body: %s", resp.StatusCode, string(b))
	}

	return
}
