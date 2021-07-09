package yearning

import (
	"encoding/json"
)

type LoginResult struct {
	Payload LoginPayload `json:"payload"`
	Code    int64        `json:"code"`
	Text    string       `json:"text"`
}

type LoginPayload struct {
	Permissions string `json:"permissions"`
	RealName    string `json:"real_name"`
	Token       string `json:"token"`
	User        string `json:"user"`
}

// Login Login yearning
func (y *Yearning) Login(normal string) (loginRes LoginResult, err error) {
	path := "/Login"
	if normal == "ldap" {
		path = "/ldap"
	}

	params := make(map[string]interface{})
	params["username"] = y.username
	params["password"] = y.password

	body, err := y.request("POST", path, nil, params)
	if err != nil {
		return
	}

	err = json.Unmarshal(body, &loginRes)
	return
}
