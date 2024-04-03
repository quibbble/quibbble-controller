package auth

import (
	"context"
	"net/http"

	"github.com/clerkinc/clerk-sdk-go/clerk"
)

type AuthInfo string

const TokenQueryParam = "token"

const (
	UserID   AuthInfo = "user_id"
	Username AuthInfo = "username"
)

const UsernameClaim = "usn"

type Auth struct {
	client clerk.Client
}

func NewAuth(pkey string) (*Auth, error) {
	client, err := clerk.NewClient(pkey)
	if err != nil {
		return nil, err
	}
	return &Auth{
		client: client,
	}, nil
}

func (a *Auth) Authenticate(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		token := r.URL.Query().Get(TokenQueryParam)
		claims, err := a.client.DecodeToken(token)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(http.StatusText(http.StatusUnauthorized)))
			return
		}

		userID := claims.Subject
		username := claims.Extra[UsernameClaim].(string)

		ctx := context.WithValue(r.Context(), UserID, userID)
		ctx = context.WithValue(ctx, Username, username)
		h.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (a *Auth) Retrieve(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		token := r.URL.Query().Get(TokenQueryParam)
		claims, err := a.client.DecodeToken(token)
		if err != nil {
			h.ServeHTTP(w, r)
			return
		}

		userID := claims.Subject
		username := claims.Extra[UsernameClaim].(string)

		ctx := context.WithValue(r.Context(), UserID, userID)
		ctx = context.WithValue(ctx, Username, username)
		h.ServeHTTP(w, r.WithContext(ctx))
	})
}
