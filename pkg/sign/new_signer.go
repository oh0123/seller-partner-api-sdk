package sign

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

const (
	SIGNED_ACCESS_TOKEN_HEADER_NAME = "x-amz-access-token"
	REFRESH_GRANT_TYPE              = "refresh_token"
	ROTATE_GRANT_TYPE               = "client_credentials"
	ROTATE_SCOPE                    = "sellingpartnerapi::client_credential:rotation"
)

type LwaAuthCredentials struct {
	ClientId     string  `json:"client_id"`
	ClientSecret string  `json:"client_secret"`
	Endpoint     string  `json:"endpoint"`
	RefreshToken string  `json:"refresh_token"`
	GrantType    string  `json:"grant_type"`
	Scope        *string `json:"scope,omitempty"`
}

type Token struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
	ExpiresIn    int    `json:"expires_in"`
}

type RotationToken struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	Scope       string `json:"scope"`
	ExpiresIn   int    `json:"expires_in"`
}

type TokenError struct {
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
}

func (t *LwaAuthCredentials) GetAccessTokenFromEndpoint(ctx context.Context) (*Token, error) {
	t.GrantType = REFRESH_GRANT_TYPE
	reqBody, err := json.Marshal(t)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, t.Endpoint, io.NopCloser(bytes.NewReader(reqBody)))
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
		var tokenError TokenError
		if err := json.Unmarshal(buf.Bytes(), &tokenError); err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("error: %v", tokenError)
	}
	var token Token
	if err := json.Unmarshal(buf.Bytes(), &token); err != nil {
		return nil, err
	}
	return &token, nil
}

func (t *LwaAuthCredentials) GetRotateAccessTokenFromEndpoint(ctx context.Context) (*RotationToken, error) {
	data := url.Values{}
	data.Set("grant_type", ROTATE_GRANT_TYPE)
	data.Set("scope", ROTATE_SCOPE)
	data.Set("client_id", t.ClientId)
	data.Set("client_secret", t.ClientSecret)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, t.Endpoint, io.NopCloser(strings.NewReader(data.Encode())))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
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
		var tokenError TokenError
		if err := json.Unmarshal(buf.Bytes(), &tokenError); err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("error: %v", tokenError)
	}
	var token RotationToken
	if err := json.Unmarshal(buf.Bytes(), &token); err != nil {
		return nil, err
	}
	return &token, nil
}

type LwaAuthSigner struct{}

func (s *LwaAuthSigner) Sign(r *http.Request, token string) {
	r.Header.Add(SIGNED_ACCESS_TOKEN_HEADER_NAME, token)
}
