# go-questrade-oauth2
Oauth2 implementation for Questrade personal apps

# Usage

```

import (
  "golang.org/x/oauth2"
  "github.com/arianitu/go-questrade-oauth2"
) 

conf := &questradeoauth2.Config{
  RefreshToken: "token from Questrade personal app"
  IsPractice: false,
}
	
client, apiServer, err := conf.Client(oauth2.NoContext)
if err != nil {
    fmt.Println(err)
    return
}

resp, err := client.Get(apiServer + "v1/time")

```
