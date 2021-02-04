package sdk

import (
	"github.com/vdobler/ht/cookiejar"
	"net/http"
)

const (
	UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/87.0.4280.88 Safari/537.36"
)

var (
	CookieJar  *cookiejar.Jar
	HttpClient *http.Client
)

func init() {
	CookieJar, _ = cookiejar.New(nil)
	HttpClient = &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
		Jar: CookieJar,
	}
}
