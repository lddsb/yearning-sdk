package yearning

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

const Unauthorized = Error("Unauthorized")

type Error string

func (e Error) Error() string { return string(e) }

func (Error) YearningError() {}

type Yearning struct {
	token     string
	username  string
	password  string
	host      string
	loginType string
	expireAt  int64
}

type JWTPayload struct {
	Exp  int64  `json:"exp"`
	Name string `json:"name"`
	Role string `json:"role"`
}

// NewClient new a client
func NewClient(username, password, host, loginType string) *Yearning {
	y := &Yearning{
		username:  username,
		password:  password,
		host:      host,
		loginType: loginType,
	}

	y.checkToken()
	return y
}

func (y *Yearning) checkToken() {
	autoLogin := func() {
		loginRes, err := y.Login(y.loginType)
		if err != nil {
			panic(fmt.Errorf("login error: %v", err))
		}

		y.token = loginRes.Payload.Token
		// 获取过期时间，将其设置到属性中
		y.expireAt = getExpireAt(y.token)
	}

	if y.expireAt < time.Now().Unix() {
		autoLogin()
	}
}

func getExpireAt(token string) int64 {
	jwtArr := strings.Split(token, ".")
	if len(jwtArr) < 3 {
		return 0
	}

	decodeString, err := base64.RawStdEncoding.DecodeString(jwtArr[1])
	if err != nil {
		return 0
	}

	var jwtPayload JWTPayload
	err = json.Unmarshal(decodeString, &jwtPayload)
	if err != nil {
		return 0
	}

	return jwtPayload.Exp
}

func commonRequest(method, url string, headers map[string]string, params map[string]interface{}) (b []byte, err error) {
	jsonData, err := json.Marshal(params)
	if err != nil {
		return
	}
	client := &http.Client{}
	req, err := http.NewRequest(method, url, bytes.NewReader(jsonData))
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

// request base request method
func (y *Yearning) request(method, path string, headers map[string]string, params map[string]interface{}) (b []byte, err error) {
	// 判断token是否已过期，如果已过期则直接抛出异常或进行登录操作
	y.checkToken()
	if headers == nil {
		headers = make(map[string]string)
	}
	headers["Authorization"] = "Bearer " + y.token
	return commonRequest(method, y.host+path, headers, params)
}
