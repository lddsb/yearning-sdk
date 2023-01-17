package yearning

import "encoding/json"

type QueryStatusResult struct {
	Payload QueryStatusPayload `json:"payload"`
	Code    int64              `json:"code"`
	Text    string             `json:"text"`
}

type QueryStatusPayload struct {
	Export bool   `json:"export"`
	Idc    string `json:"idc"`
	Status int64  `json:"status"`
}

// QueryStatus the query status
func (y *Yearning) QueryStatus() (statusRes QueryStatusResult, err error) {
	path := "/api/v2/query/status"
	b, err := y.request("PUT", path, nil, nil)
	if err != nil {
		return
	}

	err = json.Unmarshal(b, &statusRes)
	return
}

type QueryReferResult struct {
	Payload interface{} `json:"payload"`
	Code    int64       `json:"code"`
	Text    string      `json:"text"`
}

// QueryRefer submit a reference
func (y *Yearning) QueryRefer(assigned, idc, reason string) (referRes QueryReferResult, err error) {
	path := "/api/v2/query/refer"

	params := make(map[string]interface{})
	params["assigned"] = assigned
	params["export"] = 0
	params["idc"] = idc
	params["text"] = reason
	b, err := y.request("POST", path, nil, params)
	if err != nil {
		return
	}

	err = json.Unmarshal(b, &referRes)
	return
}

type QueryFetchBaseResult struct {
	Payload FetchBasePayload `json:"payload"`
	Code    int64            `json:"code"`
	Text    string           `json:"text"`
}

type FetchBasePayload struct {
	Highlight []FetchBaseHighlight `json:"highlight"`
	Idc       string               `json:"idc"`
	Info      []FetchBaseInfo      `json:"info"`
	Sign      interface{}          `json:"sign"`
	Status    int64                `json:"status"`
}

type FetchBaseHighlight struct {
	Meta string `json:"meta"`
	Vl   string `json:"vl"`
}

type FetchBaseInfo struct {
	Children []FetchBaseInfoChild `json:"children"`
	Expand   string               `json:"expand"`
	Title    string               `json:"title"`
}

type FetchBaseInfoChild struct {
	Children []FetchBaseChildChild `json:"children"`
	Title    string                `json:"title"`
}

type FetchBaseChildChild struct {
}

// QueryFetchBase get the databases from a source
func (y *Yearning) QueryFetchBase(source string) (fetchBase QueryFetchBaseResult, err error) {
	path := "/api/v2/query/fetch_base"
	params := make(map[string]interface{})
	params["source"] = source

	b, err := y.request("PUT", path, nil, params)
	if err != nil {
		return
	}

	err = json.Unmarshal(b, &fetchBase)
	return
}


type QueryResults struct {
	Payload QueryResultPayload `json:"payload"`
	Code    int64   `json:"code"`
	Text    string  `json:"text"`
}

type QueryResultPayload struct {
	Data   []map[string]string `json:"data"`
	Status bool               `json:"status"`
	Time   int64              `json:"time"`
	Title  []QueryResultTitle `json:"title"`
	Total  int64              `json:"total"`
}

type QueryResultTitle struct {
	Fixed *string `json:"fixed,omitempty"`
	Key   string  `json:"key"`
	Title string  `json:"title"`
	Width string  `json:"width"`
}

// QueryResults get sql results
func (y *Yearning) QueryResults(source, database, sql string) (qRes QueryResults, err error) {
	path := "/api/v2/query/results"
	params := make(map[string]interface{})
	params["data_base"] = database
	params["source"] = source
	params["sql"] = sql

	b, err := y.request("POST", path, nil, params)
	if err != nil {
		return
	}

	err = json.Unmarshal(b, &qRes)
	return
}