package valorant

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/url"

	tls "github.com/refraction-networking/utls"
)

type Client struct {
	httpClient *http.Client
	authCookie string
}

type UriTokens struct {
	AccessToken string
	IdToken     string
	ExpiresIn   int
}

type LoginResponseBody struct {
	Type     string `json:"type"`
	Error    string `json:"error"`
	Response struct {
		Mode       string `json:"mode"`
		Parameters struct {
			Uri string `json:"uri"`
		} `json:"parameters"`
	} `json:"response"`
	Country string `json:"country"`
}

var (
	RiotUserAgent = "RiotClient/63.0.9.4909983.4789131 rso-auth (Windows;10;;Professional, x64)"
	tlsConfig     = &tls.Config{
		MaxVersion: tls.VersionTLS13,
		MinVersion: tls.VersionTLS13,
		CipherSuites: []uint16{
			tls.TLS_AES_128_GCM_SHA256,
			tls.TLS_AES_256_GCM_SHA384,
			tls.TLS_CHACHA20_POLY1305_SHA256,
		},
	}
)

func New(proxy *url.URL) *Client {
	transport := &http.Transport{DialTLS: dialTls}
	if proxy != nil {
		transport.Proxy = http.ProxyURL(proxy)
	}

	return &Client{httpClient: &http.Client{Transport: transport}}
}

func (c *Client) Authorize(username, password string) (*UriTokens, error) {
	cookie, err := c.getPreAuth()
	if err != nil {
		return nil, err
	}

	bodyMap := map[string]any{"password": password, "type": "auth", "username": username}
	body, err := json.Marshal(bodyMap)
	if err != nil {
		return nil, err
	}

	req, err := createNewRequest("PUT", "https://auth.riotgames.com/api/v1/authorization", bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	req.Header.Add("Cookie", cookie)

	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	loginBody := new(LoginResponseBody)
	if err = json.NewDecoder(res.Body).Decode(&loginBody); err != nil {
		return nil, err
	}

	c.authCookie = parseAuthCookie(res.Header["Set-Cookie"])

	if loginBody.Type == "response" {
		return parseUriTokens(loginBody.Response.Parameters.Uri)
	} else if loginBody.Type == "auth" {
		if _, ok := ResponseErrors[loginBody.Error]; ok {
			return nil, ResponseErrors[loginBody.Error]
		}
		return nil, ErrorRiotUnknownErrorType
	} else if loginBody.Type == "multifactor" {
		return nil, ErrorRiotMultifactor
	} else {
		return nil, ErrorRiotUnknownResponseType
	}
}

func (c *Client) SubmitTwoFactor(code string) (*UriTokens, error) {
	bodyMap := map[string]any{"type": "multifactor", "code": code, "rememberDevice": true}
	body, err := json.Marshal(bodyMap)
	if err != nil {
		return nil, err
	}

	req, err := createNewRequest("PUT", "https://auth.riotgames.com/api/v1/authorization", bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	req.Header.Add("Cookie", c.authCookie)

	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	loginBody := new(LoginResponseBody)
	if err = json.NewDecoder(res.Body).Decode(&loginBody); err != nil {
		return nil, err
	}

	if loginBody.Type == "response" {
		return parseUriTokens(loginBody.Response.Parameters.Uri)
	} else if loginBody.Type == "auth" {
		if _, ok := ResponseErrors[loginBody.Error]; ok {
			return nil, ResponseErrors[loginBody.Error]
		}
		return nil, ErrorRiotUnknownErrorType
	} else if loginBody.Type == "multifactor" {
		return nil, ErrorRiotMultifactor
	} else {
		return nil, ErrorRiotUnknownResponseType
	}
}

func (c *Client) getPreAuth() (string, error) {
	nonce, err := GenerateNonce()
	if err != nil {
		return "", err
	}

	bodyMap := map[string]any{
		"acr_values": "", "claims": "",
		"client_id": "riot-client", "code_challenge": "",
		"code_challenge_method": "", "nonce": nonce,
		"redirect_uri":  "http://localhost/redirect",
		"scope":         "openid link ban lol_region account",
		"response_type": "token id_token",
	}

	body, err := json.Marshal(bodyMap)
	if err != nil {
		return "", err
	}

	req, err := createNewRequest("POST", "https://auth.riotgames.com/api/v1/authorization", bytes.NewBuffer(body))
	if err != nil {
		return "", err
	}

	res, err := c.httpClient.Do(req)
	if err != nil {
		return "", err
	}

	defer res.Body.Close()

	return parseCookies(res.Header["Set-Cookie"], "asid")
}
