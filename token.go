package main

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type User struct {
	Name string
}

type Token struct {
	Data   string    `json:"data"`
	Expiry time.Time `json:"expires_at"`
}

type AuthCache map[string]Token

var authCache AuthCache = AuthCache{}

type TokenService struct {
}

func (ts *TokenService) CreateToken(rw http.ResponseWriter, r *http.Request) {
	_, span := NewSpan(r.Context(), "create-token")
	defer span.End()

	rw.Header().Set("Content-Type", "application/json")

	data := make([]byte, 64)
	_, err := rand.Read(data)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	raw_token := base64.StdEncoding.EncodeToString(data)
	token := Token{
		Data:   raw_token,
		Expiry: time.Now().Add(time.Hour * 1),
	}

	authCache[raw_token] = token

	json.NewEncoder(rw).Encode(&token)
}

func AuthMiddleWare(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		ctx, span := NewSpan(r.Context(), "auth-middleware")
		defer span.End()

		bearer_token := r.Header.Get("Authorization")
		if bearer_token == "" {
			rw.WriteHeader(http.StatusForbidden)
			return
		}

		token := bearer_token[len("Bearer "):]
		if _, ok := authCache[token]; !ok {
			rw.WriteHeader(http.StatusForbidden)
			return
		}

		if authCache[token].Expiry.Before(time.Now()) {
			rw.WriteHeader(http.StatusForbidden)
			return
		}

		next.ServeHTTP(rw, r.Clone(ctx))
	})
}

type ComputeProvisioningService struct {
}

func (cps *ComputeProvisioningService) Provision(rw http.ResponseWriter, r *http.Request) {
	_, span := NewSpan(r.Context(), "provision-compute")
	defer span.End()

	fmt.Fprintln(rw, "Provisioning!")
}
