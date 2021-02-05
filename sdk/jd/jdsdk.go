package jd

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/yddeng/dnet/dhttp"
	"github.com/yddeng/seckill/sdk"
	"github.com/yddeng/seckill/util"
	"log"
	"math/rand"
	"net/url"
	"strconv"
	"strings"
	"time"
)

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
	req.DoEnd()
}

// 登陆二维码
func QrLoginImage(filename string) string {
	LoginPage()
	req, err := dhttp.Get(dhttp.BuildURLParams("https://qr.m.jd.com/show", url.Values{
		"appid": {"133"}, "size": {"300"}, "t": {genTime()},
	}))
	if err != nil {
		log.Println("QrLoginImage1", err.Error())
		return ""
	}
	req.Client = sdk.HttpClient
	req.SetHeader("User-Agent", sdk.UserAgent)
	req.SetHeader("Referer", "https://passport.jd.com/new/login.aspx")
	if err = req.ToFile(filename); err != nil {
		log.Println("QrLoginImage2", err.Error())
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
	return wlfstkSmdl
}

// 登陆扫码检测
func QrcodeTicket(token string) string {
	req, err := dhttp.Get(dhttp.BuildURLParams("https://qr.m.jd.com/check", url.Values{
		"appid": {"133"}, "callback": {genCallback()}, "token": {token}, "_": {genTime()},
	}))
	if err != nil {
		log.Println("QrcodeTicket1", err.Error())
		return ""
	}
	req.Client = sdk.HttpClient
	req.SetHeader("User-Agent", sdk.UserAgent)
	req.SetHeader("Referer", "https://passport.jd.com/new/login.aspx")

	type Ret struct {
		Code   int
		Msg    string
		Ticket string
	}

	var r Ret
	if body, err := req.ToString(); err != nil {
		log.Println("QrcodeTicket2", err.Error())
		return ""
	} else if err = json.Unmarshal([]byte(getCallbackStr(body)), &r); err != nil {
		log.Println("QrcodeTicket3", err.Error())
		return ""
	}

	if r.Code != 200 {
		log.Printf("Code: %d, Message: %s", r.Code, r.Msg)
		return ""
	}
	return r.Ticket
}

// 检查
func ValidQRTicket(ticket string) bool {
	req, err := dhttp.Get(dhttp.BuildURLParams("https://passport.jd.com/uc/qrCodeTicketValidation", url.Values{
		"t": {ticket},
	}))
	if err != nil {
		log.Println("ValidQRTicket1", err.Error())
		return false
	}
	req.Client = sdk.HttpClient
	req.SetHeader("User-Agent", sdk.UserAgent)
	req.SetHeader("Referer", "https://passport.jd.com/uc/login?ltype=logout")

	type Ret struct {
		ReturnCode int
		Url        string
	}
	var r Ret
	if err = req.ToJSON(&r); err != nil {
		log.Println("ValidQRTicket2", err.Error())
		return false
	}

	return r.ReturnCode == 0
}

// 获取用户数据
func GetUserNickname() string {
	req, err := dhttp.Get(dhttp.BuildURLParams("https://passport.jd.com/user/petName/getUserInfoForMiniJd.action", url.Values{
		"callback": {genCallback()}, "_": {genTime()},
	}))
	if err != nil {
		log.Println("GetUserInfo1", err.Error())
		return ""
	}
	req.Client = sdk.HttpClient
	req.SetHeader("User-Agent", sdk.UserAgent)
	req.SetHeader("Referer", "https://order.jd.com/center/list.action")

	type Ret struct {
		NickName    string
		RealName    string
		UserScoreVO struct {
			FinanceScore int // 金融分数
			TotalScore   int // 分数
		}
	}

	var r Ret
	if body, err := req.ToString(); err != nil {
		log.Println("GetUserInfo2", err.Error())
		return ""
	} else if b, err := util.GbkToUtf8([]byte(getCallbackStr(body))); err != nil {
		log.Println("GetUserInfo3", err.Error())
		return ""
	} else if err = json.Unmarshal(b, &r); err != nil {
		log.Println("GetUserInfo4", err.Error())
		return ""
	}

	if r.NickName != "" {
		log.Printf("UserInfo<NickName：%s, RealName：%s, FinanceScore：%d>\n", r.NickName, r.RealName, r.UserScoreVO.FinanceScore)
	}
	return r.NickName
}

// 验证cookie
func ValidCookie() bool {
	req, err := dhttp.Get(dhttp.BuildURLParams("https://order.jd.com/center/list.action", url.Values{
		"rid": {genTime()},
	}))
	if err != nil {
		log.Println("ValidCookie1", err.Error())
		return false
	}
	req.Client = sdk.HttpClient
	req.SetHeader("User-Agent", sdk.UserAgent)

	resp, err := req.Do()
	if err == nil && resp.StatusCode == 200 {
		defer resp.Body.Close()
		return true
	}
	return false
}

/* *********** * ************* */

// 获取服务器时间
func GetServerTime() int64 {
	req, err := dhttp.Get("https://a.jd.com//ajax/queryServerData.html")
	if err != nil {
		log.Println("GetServerTime1", err.Error())
		return 0
	}
	type Ret struct {
		ServerTime int64 //1609878734768
	}
	var r Ret
	if err = req.ToJSON(&r); err != nil {
		log.Println("GetServerTime2", err.Error())
		return 0
	}
	return r.ServerTime
}

/* ************* 商品相关 ******************** */

//获取商品信息
func GetProductInfo(skuId string) string {
	req, err := dhttp.Get(fmt.Sprintf("https://item.jd.com/%s.html", skuId))
	if err != nil {
		return ""
	}
	req.Client = sdk.HttpClient
	req.SetHeader("User-Agent", sdk.UserAgent)

	if body, err := req.ToString(); err != nil {
		return ""
	} else {
		html := strings.NewReader(body)
		doc, _ := goquery.NewDocumentFromReader(html)
		return strings.TrimSpace(doc.Find(".sku-name").Text())
	}
}

// 加入购物车
func InitCart(pid, pcount string) bool {
	// https://cart.jd.com/addToCart.html?rcd=1&pid=100015185396&pc=1&eb=1&rid=1612494295861&em=
	req, err := dhttp.Get(dhttp.BuildURLParams("https://cart.jd.com/gate.action", url.Values{
		"pid": {pid}, "pcount": {pcount}, "ptype": {"1"},
		//"rcd": {"1"}, "pid": {pid}, "pc": {"1"}, "eb": {"1"}, "rid": {"1612494295861"}, "em": {""},
	}))
	if err != nil {
		log.Println("GetSeckillInitInfo", err.Error())
		return false
	}
	req.Client = sdk.HttpClient
	req.SetHeader("User-Agent", sdk.UserAgent)
	req.SetHeader("Referer", fmt.Sprintf("https://item.jd.com/%s.html", pid))
	req.SetHeader("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	req.SetHeader("Host", "cart.jd.com")

	str, err := req.ToString()
	if err != nil {
		return false
	}
	log.Println(str)
	return true
}

func CartIndex() {
	req, err := dhttp.Get("https://cart.jd.com/cart_index/")
	if err != nil {
		return
	}
	req.Client = sdk.HttpClient
	req.SetHeader("User-Agent", sdk.UserAgent)
	req.SetHeader("Host", "cart.jd.com")
	req.SetHeader("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")

	if body, err := req.ToString(); err != nil {
		return
	} else {
		fmt.Println(body)
	}
}

// 获取秒杀商品信息
func GetSeckillInitInfo(skuId, skuNum string) (*InitData, error) {
	req, err := dhttp.NewRequest("https://marathon.jd.com/seckillnew/orderService/pc/init.action", "POST")
	if err != nil {
		log.Println("GetSeckillInitInfo1", err.Error())
		return nil, err
	}
	req.Client = sdk.HttpClient
	req.SetHeader("User-Agent", sdk.UserAgent)
	req.SetHeader("Host", "marathon.jd.com")

	req.WriteParam(url.Values{"sku": {skuId}, "num": {skuNum}, "isModifyAddress": {"false"}})

	var initData InitData
	if err = req.ToJSON(&initData); err != nil {
		log.Println("GetSeckillInitInfo2", err.Error())
		return nil, err
	} else if len(initData.AddressList) == 0 {
		return nil, errors.New("初始化秒杀信息失败, AddressList为空")
	}
	return &initData, nil
}

// 获取秒杀链接
func GetKillUrl(skuId string) string {
	req, err := dhttp.Get(dhttp.BuildURLParams("https://itemko.jd.com/itemShowBtn", url.Values{
		"skuId": {skuId}, "callback": {genCallback()}, "from": {"pc"}, "_": {genTime()},
	}))
	if err != nil {
		log.Println("GetKillUrl1", err.Error())
		return ""
	}
	req.Client = sdk.HttpClient
	req.SetHeader("User-Agent", sdk.UserAgent)
	req.SetHeader("Host", "itemko.jd.com")
	req.SetHeader("Referer", fmt.Sprintf("https://item.jd.com/%s.html", skuId))

	type Ret struct {
		Url string
	}

	var r Ret
	if body, err := req.ToString(); err != nil {
		log.Println("GetKillUrl2", err.Error())
		return ""
	} else if err = json.Unmarshal([]byte(getCallbackStr(body)), &r); err != nil {
		log.Println("GetKillUrl3", err.Error())
		return ""
	}

	if r.Url == "" {
		return ""
	}

	//https://divide.jd.com/user_routing?skuId=8654289&sn=c3f4ececd8461f0e4d7267e96a91e0e0&from=pc
	url := strings.ReplaceAll(r.Url, "divide", "marathon")
	//https://marathon.jd.com/captcha.html?skuId=8654289&sn=c3f4ececd8461f0e4d7267e96a91e0e0&from=pc
	url = strings.ReplaceAll(url, "user_routing", "captcha.html")
	return "https" + url
}

// 请求秒杀链接
func RequestKillUrl(skuId, killUrl string) bool {
	req, err := dhttp.Get(killUrl)
	if err != nil {
		log.Println("RequestKillUrl1", err.Error())
		return false
	}
	req.Client = sdk.HttpClient
	req.SetHeader("User-Agent", sdk.UserAgent)
	req.SetHeader("Host", "marathon.jd.com")
	req.SetHeader("Referer", fmt.Sprintf("https://item.jd.com/%s.html", skuId))

	resp, err := req.Do()
	if err == nil && resp.StatusCode == 200 {
		defer resp.Body.Close()
		return true
	}
	return false
}

// 访问抢购订单结算页面
func SeckillPage(skuId, skuNum string) bool {
	req, err := dhttp.Get(dhttp.BuildURLParams("https://marathon.jd.com/seckill/seckill.action", url.Values{
		"sku": {skuId}, "num": {skuNum}, "rid": {genTime()},
	}))
	if err != nil {
		log.Println("SeckillPage1", err.Error())
		return false
	}
	req.Client = sdk.HttpClient
	req.SetHeader("User-Agent", sdk.UserAgent)
	req.SetHeader("Host", "marathon.jd.com")
	req.SetHeader("Referer", fmt.Sprintf("https://item.jd.com/%s.html", skuId))

	resp, err := req.Do()
	if err == nil && resp.StatusCode == 200 {
		defer resp.Body.Close()
		return true
	}
	return false
}

// 提交订单
func SubmitSeckillOrder(eid, fp, skuId, skuNum, pwd string, initData *InitData) bool {
	req, err := dhttp.NewRequest(dhttp.BuildURLParams("https://marathon.jd.com/seckill/seckill.action", url.Values{"skuId": {skuId}}), "POST")
	if err != nil {
		log.Println("SubmitSeckillOrder1", err.Error())
		return false
	}
	req.Client = sdk.HttpClient
	req.SetHeader("User-Agent", sdk.UserAgent)
	req.SetHeader("Host", "marathon.jd.com")
	req.SetHeader("Referer", fmt.Sprintf("https://marathon.jd.com/seckill/seckill.action?skuId=%s&num=%s&rid=%s", skuId, skuNum, genTime()))

	params := url.Values{
		"eid":      {eid},
		"fp":       {fp},
		"skuId":    {skuId},
		"num":      {skuNum},
		"password": {pwd},

		"addressId":     {strconv.Itoa(initData.AddressList[0].Id)},
		"name":          {initData.AddressList[0].Name},
		"provinceId":    {strconv.Itoa(initData.AddressList[0].ProvinceId)},
		"cityId":        {strconv.Itoa(initData.AddressList[0].CityId)},
		"countyId":      {strconv.Itoa(initData.AddressList[0].CountyId)},
		"townId":        {strconv.Itoa(initData.AddressList[0].TownId)},
		"addressDetail": {initData.AddressList[0].AddressDetail},
		"mobile":        {initData.AddressList[0].Mobile},
		"mobileKey":     {initData.AddressList[0].MobileKey},
		"email":         {initData.AddressList[0].Email},

		"token": {initData.Token},

		"yuShou": {"true"},

		"isModifyAddress":    {"false"},
		"postCode":           {""},
		"invoiceCompanyName": {""},
		"invoiceTaxpayerNO":  {""},
		"invoiceEmail":       {""},
		"codTimeType":        {"3"},
		"paymentType":        {"4"},
		"areaCode":           {""},
		"overseas":           {"0"},
		"phone":              {""},
		"pru":                {""},
	}

	if initData.InvoiceInfo != nil {
		params.Set("invoice", "true")
		params.Set("invoicePhone", initData.InvoiceInfo.InvoicePhone)
		params.Set("invoicePhoneKey", initData.InvoiceInfo.InvoicePhoneKey)
		params.Set("invoiceContent", strconv.Itoa(initData.InvoiceInfo.InvoiceContentType))
		params.Set("invoiceTitle", strconv.Itoa(initData.InvoiceInfo.InvoiceTitle))
	} else {
		params.Set("invoice", "false")
	}

	req.WriteParam(params)

	type Ret struct {
		Success      bool
		ErrorMessage string
		OrderId      int
		ResultCode   int
		TotalMoney   float32
		PcUrl        string
	}

	var r Ret
	if err = req.ToJSON(&r); err != nil {
		str, _ := req.ToString()
		log.Println("SubmitSeckillOrder2", str, err.Error())
		return false
	}
	if r.Success {
		log.Println(fmt.Sprintf("抢购成功，订单号:%d, 总价:%f, 电脑端付款链接:%s", r.OrderId, r.TotalMoney, r.PcUrl))
	}
	return r.Success
}
