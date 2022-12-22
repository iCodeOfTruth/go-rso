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
    "github.com/iCodeOfTruth/go-valorant"
    "net/url"
)

func main() {
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
