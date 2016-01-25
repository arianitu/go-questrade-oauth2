package questradeoauth2

import (
	"encoding/json"
	"errors"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"net/http"
	"net/url"
	"time"
)

const (
	authURL         = "https://login.questrade.com/oauth2/token?grant_type=refresh_token&refresh_token="
	practiceAuthURL = "https://practicelogin.questrade.com/oauth2/token?grant_type=refresh_token&refresh_token="
)

type Config struct {
	RefreshToken string
	IsPractice   bool
}

// TokenSource returns a Questrade TokenSource using the configuration
// in c and the HTTP client from the provided context.
func (c *Config) TokenSource(ctx context.Context) oauth2.TokenSource {
	return oauth2.ReuseTokenSource(nil, questradeSource{ctx, c})
}

// Client returns an HTTP client wrapping the context's
// HTTP transport and adding Authorization headers with tokens
// obtained from c.
//
// The returned client and its Transport should not be modified.
func (c *Config) Client(ctx context.Context) (client *http.Client, apiServer string, err error) {
	source := c.TokenSource(ctx)
	token, err := source.Token()
	if err != nil {
		return nil, "", err
	}
	apiServer, ok := token.Extra("ApiServer").(string)
	if !ok {
		return nil, "", errors.New("ApiServer was not filled up properly")
	}

	return oauth2.NewClient(ctx, source), apiServer, nil
}

// questradeSource is a source that gets an access token based on an existing refresh token
// that comes from a Questrade personal app
type questradeSource struct {
	ctx  context.Context
	conf *Config
}

type authResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int64  `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	ApiServer    string `json:"api_server"`
}

func (qs questradeSource) Token() (*oauth2.Token, error) {
	qc := oauth2.NewClient(qs.ctx, nil)

	var apiUrl string
	if qs.conf.IsPractice {
		apiUrl = practiceAuthURL
	} else {
		apiUrl = authURL
	}

	resp, err := qc.Get(apiUrl + qs.conf.RefreshToken)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode == http.StatusBadRequest {
		return nil, errors.New("Invalid Refresh Token")
	}

	defer resp.Body.Close()
	decoder := json.NewDecoder(resp.Body)

	var authResp authResponse
	err = decoder.Decode(&authResp)
	if err != nil {
		return nil, err
	}

	token := &oauth2.Token{
		AccessToken: authResp.AccessToken,
		TokenType:   authResp.TokenType,
	}
	extra := url.Values{}
	extra.Add("ApiServer", authResp.ApiServer)

	token = token.WithExtra(extra)
	if secs := authResp.ExpiresIn; secs > 0 {
		token.Expiry = time.Now().Add(time.Duration(secs) * time.Second)
	}

	qs.conf.RefreshToken = authResp.RefreshToken
	return token, nil
}
