package questradeoauth2

import (
	"encoding/json"
	"fmt"
	"golang.org/x/oauth2"
	"os"
	"testing"
	"time"
)

func TestAuthentication(t *testing.T) {
	conf := &Config{
		RefreshToken: os.Getenv("REFRESH_TOKEN"),
		IsPractice:   false,
	}

	client, apiServer, err := conf.Client(oauth2.NoContext)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(apiServer)

	resp, err := client.Get(apiServer + "v1/time")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	decoder := json.NewDecoder(resp.Body)

	var timeResponse struct {
		Time string `json:"time"`
	}

	err = decoder.Decode(&timeResponse)
	if err != nil {
		t.Fatal(err)
	}

	_, err = time.Parse(time.RFC3339, timeResponse.Time)
	if err != nil {
		t.Fatal(err)
	}
}
