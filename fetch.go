package yearning

import "encoding/json"

type FetchIdcResult struct {
	Payload []string `json:"payload"`
	Code    int64    `json:"code"`
	Text    string   `json:"text"`
}

// FetchIDC get idc list
func (y *Yearning) FetchIDC() (idcRes FetchIdcResult, err error) {
	path := "/api/v2/fetch/idc"
	b, err := y.request("GET", path, nil, nil)
	if err != nil {
		return
	}

	err = json.Unmarshal(b, &idcRes)
	return
}

type FetchSourceResult struct {
	Payload FetchSourcePayload `json:"payload"`
	Code    int64   `json:"code"`
	Text    string  `json:"text"`
}

type FetchSourcePayload struct {
	Assigned []string `json:"assigned"`
	Source   []string `json:"source"`
}

// FetchSource get sources from a IDC
func (y *Yearning) FetchSource(idc, tp string) (sourceRes FetchSourceResult, err error) {
	path := "/api/v2/fetch/source"
	// TODO params
	params := make(map[string]interface{})
	params["idc"] = idc
	params["tp"] = tp

	b, err := y.request("GET", path, nil, params)
	if err != nil {
		return
	}

	err = json.Unmarshal(b, &sourceRes)
	return
}

type FetchDataBaseResult struct {
	Payload FetchDatabasePayload `json:"payload"`
	Code    int64   `json:"code"`
	Text    string  `json:"text"`
}

type FetchDatabasePayload struct {
	Highlight []FetchDatabaseHighlight `json:"highlight"`
	Results   []string    `json:"results"`
}

type FetchDatabaseHighlight struct {
	Meta string `json:"meta"`
	Vl   string `json:"vl"`
}

// FetchDatabases get databases from a source
func (y *Yearning) FetchDatabases(source string) (baseRes FetchDataBaseResult, err error) {
	path := "/api/v2/fetch/base"

	params := make(map[string]interface{})
	params["source"] = source

	b, err := y.request("GET", path, nil, params)
	if err != nil {
		return
	}

	err = json.Unmarshal(b, &baseRes)
	return
}

type FetchTablesResult struct {
	Payload FetchTablesPayload `json:"payload"`
	Code    int64   `json:"code"`
	Text    string  `json:"text"`
}

type FetchTablesPayload struct {
	Highlight []FetchTablesHighlight `json:"highlight"`
	Table     []string    `json:"table"`
}

type FetchTablesHighlight struct {
	Meta string `json:"meta"`
	Vl   string `json:"vl"`
}


type TableQuery struct {
	IDC string `json:"idc"`
	Source string `json:"source"`
	Database string `json:"data_base"`
	Table string `json:"table"`
	Reason string `json:"reason"`
	Delay string `json:"delay"`
	Assigned string `json:"assigned"`
	Backup int `json:"backup"`
	Export int `json:"export"`
	Tp int `json:"tp"`
}

// FetchTables fetch tables
func (y *Yearning) FetchTables(tQ *TableQuery) (tRes FetchTablesResult, err error) {
	path := "/api/v2/fetch/table"
	params := make(map[string]interface{})
	marshal, err := json.Marshal(tQ)
	if err != nil {
		return
	}
	err = json.Unmarshal(marshal, &params)
	if err != nil {
		return
	}

	b, err := y.request("GET", path, nil, params)
	if err != nil {
		return
	}

	var a interface{}
	err = json.Unmarshal(b, &a)
	return
}

