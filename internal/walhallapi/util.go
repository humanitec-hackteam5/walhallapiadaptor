package walhallapi

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	jwt "github.com/dgrijalva/jwt-go"
)

type APIState struct {
	claims    WalhallClaims
	jwt       string
	apiPrefix string
	doer      Doer
	cache     map[string]interface{}
}

type Doer interface {
	Do(req *http.Request) (*http.Response, error)
}

type WalhallClaims struct {
	jwt.StandardClaims
	UserUUID string   `json:"user_uuid,omitempty"`
	OrgUUIDs []string `json:"organization_uuids,omitempty"`
	Username string   `json:"username,omitempty"`
	Scope    string   `json:"scope,omitempty"`
}

func New(apiPrefix, jwt string, doer Doer) (*APIState, error) {
	// There is some funkiness with when the JWT string gets garbage collected, so use replace to guarantee a copy
	jwt = strings.Replace(jwt, "JWT ", "", 1)
	claims, err := claimsFromJWT(jwt)
	if err != nil {
		return nil, err
	}
	return &APIState{
		claims:    claims,
		jwt:       "JWT " + jwt,
		apiPrefix: apiPrefix,
		doer:      doer,
		cache:     make(map[string]interface{}),
	}, nil
}

func claimsFromJWT(JWT string) (WalhallClaims, error) {
	parser := jwt.Parser{nil, false, true}
	var claims WalhallClaims
	_, _, err := (&parser).ParseUnverified(JWT, &claims)
	if err != nil {
		return WalhallClaims{}, fmt.Errorf("exracting walhall claims from JWT: %v", err)
	}
	return claims, nil
}

func (a *APIState) makeRequest(method string, url string, body io.Reader) (*http.Response, error) {
	var req *http.Request
	var err error
	// We need to handle typed and non-typed nils (see: https://golang.org/doc/faq#nil_error)
	if body == nil {
		req, err = http.NewRequest(method, a.apiPrefix+url, nil)
	} else {
		req, err = http.NewRequest(method, a.apiPrefix+url, body)
	}
	if err != nil {
		return nil, fmt.Errorf("make request: %w", err)
	}
	req.Header["authorization"] = []string{a.jwt}
	req.Header["Content-Type"] = []string{"application/json"}
	log.Printf("[walhallapi] %s %s%s", method, a.apiPrefix, url)
	return a.doer.Do(req)
}
