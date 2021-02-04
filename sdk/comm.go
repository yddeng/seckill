package sdk

import (
	"encoding/json"
	"github.com/vdobler/ht/cookiejar"
	"io/ioutil"
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

func SaveCookie(filename string) {
	cookies := make(map[string]cookiejar.Entry)
	for _, tld := range CookieJar.ETLDsPlus1(nil) {
		for _, cookie := range CookieJar.Entries(tld, nil) {
			id := cookie.ID()
			cookies[id] = cookie
		}
	}
	data, err := json.MarshalIndent(cookies, "", "    ")
	if err == nil {
		_ = ioutil.WriteFile(filename, data, 0666)
	}
}

func LoadCookie(filename string) bool {
	if filename == "" {
		return false
	}
	buf, err := ioutil.ReadFile(filename)
	if err != nil {
		return false
	}

	cookies := make(map[string]cookiejar.Entry)
	err = json.Unmarshal(buf, &cookies)
	if err != nil {
		return false
	}
	cs := make([]cookiejar.Entry, 0, len(cookies))
	for _, c := range cookies {
		cs = append(cs, c)
	}

	jar, _ := cookiejar.New(nil)
	jar.LoadEntries(cs)

	CookieJar = jar
	HttpClient.Jar = CookieJar
	return true
}
