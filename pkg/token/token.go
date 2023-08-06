package token

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const (
	GrantType = "refresh_token"
)

type Token struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
	ExpiresIn    int    `json:"expires_in"`
}

type TokenError struct {
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
}

type Config struct {
	ClientID     string  `json:"client_id"`
	ClientSecret string  `json:"client_secret"`
	RefreshToken string  `json:"refresh_token"`
	GrantType    string  `json:"grant_type"`
	Scope        *string `json:"scope,omitempty"`
}

func (t *Config) GetAccessToken(ctx context.Context, endpoint string) (*Token, error) {
	t.GrantType = GrantType
	reqBody, err := json.Marshal(t)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, "POST", endpoint, io.NopCloser(bytes.NewReader(reqBody)))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json;charset=UTF-8")
	client := http.DefaultClient
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode >= http.StatusBadRequest {
		tokenError := new(TokenError)
		if err := json.Unmarshal(buf.Bytes(), tokenError); err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("error: %v", tokenError)
	}
	token := new(Token)
	if err := json.Unmarshal(buf.Bytes(), token); err != nil {
		return nil, err
	}
	return token, nil
}
