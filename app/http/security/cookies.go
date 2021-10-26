package security

import (
	"github.com/gin-gonic/contrib/sessions"
	"os"
)

type CookiesService struct {
	Store sessions.CookieStore
}

func NewCookies() (*CookiesService, error) {
	store := sessions.NewCookieStore([]byte(os.Getenv("COOKIES_SECRET")))
	return &CookiesService{
		Store: store,
	}, nil
}
