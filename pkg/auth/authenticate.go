package auth

import (
	"context"
	"net/http"
	"strings"

	"github.com/clerkinc/clerk-sdk-go/clerk"
)

type AuthInfo string

const (
	UserID   AuthInfo = "UserID"
	Username AuthInfo = "Username"
)

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

		authorization := r.Header.Get("Authorization")
		s := strings.Split(authorization, " ")
		if len(s) < 2 {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(http.StatusText(http.StatusUnauthorized)))
			return
		}
		token := s[1]
		claims, err := a.client.DecodeToken(token)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(http.StatusText(http.StatusUnauthorized)))
			return
		}
		user, err := a.client.Users().Read(claims.Subject)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(http.StatusText(http.StatusUnauthorized)))
			return
		}

		ctx := context.WithValue(r.Context(), UserID, user.ID)
		ctx = context.WithValue(ctx, Username, user.Username)
		h.ServeHTTP(w, r.WithContext(ctx))
	})
}
