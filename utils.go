package valorant

import (
	"fmt"
	tls "github.com/refraction-networking/utls"
	"io"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

func parseUriTokens(uri string) (*UriTokens, error) {
	u, err := url.Parse(uri)
	if err != nil {
		return nil, err
	}

	q, err := url.ParseQuery(u.Fragment)
	if err != nil {
		return nil, err
	}

	accessToken := q.Get("access_token")
	idToken := q.Get("id_token")

	expiresIn, err := strconv.Atoi(q.Get("expires_in"))
	if err != nil {
		return nil, err
	}

	return &UriTokens{
		AccessToken: accessToken,
		IdToken:     idToken,
		ExpiresIn:   expiresIn,
	}, nil
}

func parseCookies(cookies []string, subs string) (string, error) {
	for _, cookie := range cookies {
		if strings.Contains(cookie, subs) {
			return cookie, nil
		}
	}

	return "", fmt.Errorf("could not find %s", subs)
}

func dialTls(network, addr string) (net.Conn, error) {
	netConn, err := net.Dial(network, addr)
	if err != nil {
		return nil, err
	}

	host, _, err := net.SplitHostPort(addr)
	if err != nil {
		return nil, err
	}

	config := tlsConfig.Clone()
	config.ServerName = host

	tlsConn := tls.UClient(netConn, config, tls.HelloGolang)
	if err = tlsConn.Handshake(); err != nil {
		return nil, err
	}

	return tlsConn, nil
}

func createNewRequest(method, url string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	req.Header = http.Header{
		"Cache-Control": {"no-cache"},
		"Content-Type":  {"application/json"},
		"Cookie":        {""},
		"User-Agent":    {RiotUserAgent},
	}

	return req, nil
}
