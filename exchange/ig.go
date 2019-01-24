package exchange

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

// OAuthToken goten at login
type OAuthToken struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
	TokenType    string `json:"token_type"`
	ExpiresIn    string `json:"expires_in"`
}

// Session session struct
type Session struct {
	ClientID              string     `json:"clientId"`
	AccountID             string     `json:"accountId"`
	TimezoneOffset        int        `json:"timezoneOffset"`
	LightstreamerEndpoint string     `json:"lightstreamerEndpoint"`
	Token                 OAuthToken `json:"oauthToken"`
}

// IGClient connects to IG
type IGClient struct {
	host    string
	apiKey  string
	session *Session
}

// ============================ PRIVATE ============================

func (ig *IGClient) header() map[string]string {
	contentType := "application/json; charset=UTF-8"

	header := make(map[string]string)
	header["Content-Type"] = contentType
	header["Accept"] = contentType
	header["X-IG-API-KEY"] = ig.apiKey

	if ig.session != nil {
		header["IG-ACCOUNT-ID"] = ig.session.AccountID
		header["Authorization"] = fmt.Sprintf("%s %s", ig.session.Token.TokenType, ig.session.Token.AccessToken)
	}

	return header
}

func (ig *IGClient) request(method, url string, body io.Reader, version string) (*http.Response, error) {
	client := &http.Client{}

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	for key, value := range ig.header() {
		req.Header.Set(key, value)
	}
	req.Header.Set("Version", version)

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return resp, nil
}

func (ig *IGClient) login(username, password string) *Session {

	payload := struct {
		Identifier string `json:"identifier"`
		Password   string `json:"password"`
	}{username, password}
	payloadBytes, _ := json.Marshal(payload)

	url := ig.host + "/session"
	resp, err := ig.request("POST", url, bytes.NewReader(payloadBytes), "3")

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}

	response := Session{}

	json.Unmarshal(body, &response)

	return &response
}

// ============================ PUBLIC ============================

// NewIGClient create new IG client
func NewIGClient(username, password, apiKey string) *IGClient {
	ig := new(IGClient)
	ig.host = "https://demo-api.ig.com/gateway/deal"
	ig.apiKey = apiKey
	ig.session = ig.login(username, password)

	return ig
}

// Balance struct
type Balance struct {
	Balance    float64 `json:"balance"`
	Deposit    float64 `json:"deposit"`
	ProfitLoss float64 `json:"profitLoss"`
	Available  float64 `json:"available"`
}

// Account struct
type Account struct {
	AccountID       string  `json:"accountId"`
	AccountName     string  `json:"accountName"`
	AccountAlias    string  `json:"accountAlias"`
	Status          string  `json:"status"`
	AccountType     string  `json:"accountType"`
	Preferred       bool    `json:"preferred"`
	Balance         Balance `json:"balance"`
	Currency        string  `json:"currency"`
	CanTransferFrom bool    `json:"canTransferFrom"`
	CanTransferTo   bool    `json:"canTransferTo"`
}

// Accounts get IG accounts
func (ig *IGClient) Accounts() []Account {
	url := ig.host + "/accounts"
	resp, err := ig.request("GET", url, nil, "1")

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}

	response := struct {
		Accounts []Account `json:"accounts"`
	}{}

	json.Unmarshal(body, &response)

	return response.Accounts
}
