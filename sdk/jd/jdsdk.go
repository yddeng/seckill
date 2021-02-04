package jd

import (
	"encoding/json"
	"github.com/vdobler/ht/cookiejar"
	"github.com/yddeng/dnet/dhttp"
	"github.com/yddeng/seckill/util"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/87.0.4280.88 Safari/537.36"

var (
	cookieJar  *cookiejar.Jar
	httpClient *http.Client
)

func init() {
	cookieJar, _ = cookiejar.New(nil)
	httpClient = &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
		Jar: cookieJar,
	}
}

func genTime() string {
	return strconv.Itoa(time.Now().Second() * 1000)
}

func genCallback() string {
	return "jQuery" + strconv.Itoa(int(1000000+rand.Int31n(8999999)))
}

func getCallbackStr(text string) string {
	fromIndex := strings.Index(text, "{")
	endIndex := strings.LastIndex(text, "}")
	return text[fromIndex : endIndex+1]
}

// 登录页面
func LoginPage() {
	req, _ := dhttp.Get("https://passport.jd.com/new/login.aspx")
	_, _ = req.Do()
}

// 登陆二维码
func QrLoginImage() string {
	LoginPage()
	req, err := dhttp.Get(dhttp.BuildURLParams("https://qr.m.jd.com/show", url.Values{
		"appid": {"133"}, "size": {"300"}, "t": {genTime()},
	}))
	if err != nil {
		log.Panicln("QrLoginImage", err.Error())
		return ""
	}
	req.Client = httpClient
	req.SetHeader("User-Agent", UserAgent)
	req.SetHeader("Referer", "https://passport.jd.com/new/login.aspx")
	if err = req.ToFile("./jd_qr_code.png"); err != nil {
		log.Panicln("获取二维码失败")
		return ""
	}

	cookies := req.HttpResponse().Cookies()
	wlfstkSmdl := ""
	for _, cookie := range cookies {
		if cookie.Name == "wlfstk_smdl" {
			wlfstkSmdl = cookie.Value
			break
		}
	}
	if wlfstkSmdl != "" {
		log.Println("二维码获取成功，请打开京东APP扫描")
		util.OpenImage("./jd_qr_code.png")
	}
	return wlfstkSmdl
}

// 登陆扫码检测
func QrcodeTicket(token string) string {
	req, err := dhttp.Get(dhttp.BuildURLParams("https://qr.m.jd.com/check", url.Values{
		"appid": {"133"}, "callback": {genCallback()}, "token": {token}, "_": {genTime()},
	}))
	if err != nil {
		log.Panicln("QrcodeTicket1", err.Error())
		return ""
	}
	req.Client = httpClient
	req.SetHeader("User-Agent", UserAgent)
	req.SetHeader("Referer", "https://passport.jd.com/new/login.aspx")

	type Ret struct {
		Code   int
		Msg    string
		Ticket string
	}

	var r Ret
	if body, err := req.ToString(); err != nil {
		log.Panicln("QrcodeTicket2", err.Error())
		return ""
	} else if err = json.Unmarshal([]byte(getCallbackStr(body)), &r); err != nil {
		log.Panicln("QrcodeTicket3", err.Error())
		return ""
	}

	if r.Code != 200 {
		log.Printf("Code: %d, Message: %s", r.Code, r.Msg)
		return ""
	}
	log.Println("已完成手机客户端确认")
	return r.Ticket
}

// 检查
func ValidQRTicket(ticket string) bool {
	req, err := dhttp.Get(dhttp.BuildURLParams("https://passport.jd.com/uc/qrCodeTicketValidation", url.Values{
		"t": {ticket},
	}))
	if err != nil {
		log.Panicln("ValidQRTicket1", err.Error())
		return false
	}
	req.Client = httpClient
	req.SetHeader("User-Agent", UserAgent)
	req.SetHeader("Referer", "https://passport.jd.com/uc/login?ltype=logout")

	type Ret struct {
		ReturnCode int
		Url        string
	}
	var r Ret
	if err = req.ToJSON(&r); err != nil {
		log.Panicln("ValidQRTicket1", err.Error())
		return false
	}

	return r.ReturnCode == 0
}

// 获取用户数据
func GetUserInfo() string {
	req, err := dhttp.Get(dhttp.BuildURLParams("https://passport.jd.com/user/petName/getUserInfoForMiniJd.action", url.Values{
		"callback": {genCallback()}, "_": {genTime()},
	}))
	if err != nil {
		log.Panicln("GetUserInfo1", err.Error())
		return ""
	}
	req.Client = httpClient
	req.SetHeader("User-Agent", UserAgent)
	req.SetHeader("Referer", "https://order.jd.com/center/list.action")

	type Ret struct {
		NickName string
	}

	var r Ret
	if body, err := req.ToString(); err != nil {
		log.Panicln("GetUserInfo2", err.Error())
		return ""
	} else if err = json.Unmarshal([]byte(getCallbackStr(body)), &r); err != nil {
		log.Panicln("GetUserInfo3", err.Error())
		return ""
	}

	b, _ := util.GbkToUtf8([]byte(r.NickName))
	return string(b)
}

// 验证cookie
func ValidCookie() bool {
	req, err := dhttp.Get(dhttp.BuildURLParams("https://order.jd.com/center/list.action", url.Values{
		"rid": {genTime()},
	}))
	if err != nil {
		log.Panicln("ValidCookie", err.Error())
		return false
	}
	req.Client = httpClient
	req.SetHeader("User-Agent", UserAgent)

	resp, err := req.Do()
	if err == nil && resp.StatusCode == 200 {
		return true
	}
	return false
}

/* *********** * ************* */
