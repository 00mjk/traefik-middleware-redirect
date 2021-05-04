package traefik_middleware_redirect

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strconv"
)

// Config the plugin configuration.
type Config struct {
	RedirectCode int    `json:"redirectCode,omitempty"`
	RedirectURI  string `json:"redirectUri,omitempty"`
	ExpiringTime int    `json:"expiringTime,omitempty"`
}

func CreateConfig() *Config {
	return &Config{
		RedirectCode: 301,
		RedirectURI:  "/",
		ExpiringTime: 0,
	}
}

type Redirect struct {
	next         http.Handler
	redirectCode int
	redirectURI  string
	expiringTime int
}

func New(ctx context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {

	//log.Printf("Redirect Plugin. Code: %d, URI: %s, Expired Time: %d", config.RedirectCode, config.RedirectURI, config.ExpiringTime)

	return &Redirect{
		redirectCode: config.RedirectCode,
		redirectURI:  config.RedirectURI,
		expiringTime: config.ExpiringTime,
		next:         next,
	}, nil

}

func (middleware *Redirect) redirect() http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		//recorder is used to delegate call. Its response will be used to create the correct ResponseWriter status code.
		recorder := httptest.NewRecorder()
		middleware.next.ServeHTTP(recorder, req)

		if middleware.expiringTime > 0 {
			rw.Header().Set("Cache-Control", "max-age:"+strconv.Itoa(middleware.expiringTime))
		}

		http.Redirect(
			rw,
			req,
			middleware.redirectURI,
			middleware.redirectCode,
		)

	})
}

func (middleware *Redirect) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	middleware.redirect().ServeHTTP(rw, req)
}
