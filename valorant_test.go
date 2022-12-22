package valorant

import (
	"net/url"
	"testing"
)

func TestValorantAuth(t *testing.T) {
	proxyUrl, _ := url.Parse("http://user:pass@ip:port")
	client := New(proxyUrl)

	data, err := client.Authorize("Username", "Password")
	if err != nil {
		t.Error(err)
	}

	t.Log(data)
}
