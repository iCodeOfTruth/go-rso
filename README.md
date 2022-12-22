# ValorantAuth
Riot Sign-in in Go

# Credits
Big thanks to https://github.com/fyraux/go-rso for the base of this project

# Errors
Check file `errors.go` for all errors

# Usage

```go
package valorant

import (
    "fmt"
    "net/url"

	valorant "github.com/iCodeOfTruth/go-rso"
)

func main() {
	valorant.RiotUserAgent = "RiotClient/62.0.1.4852117.4789131 rso-auth (Windows;11;;Professional, x64)" // Set your own user agent

	// Proxy support
	proxyUrl, _ := url.Parse("http://user:pass@ip:port")
	client := valorant.New(proxyUrl) // or New(nil) for no proxy

	data, err := client.Authorize("Username", "Password")
	if err != nil {
		panic(err)
	}

	fmt.Println(data.AccessToken)
}
```
