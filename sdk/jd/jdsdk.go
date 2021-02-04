package jd

import (
	"encoding/json"
	"errors"
	"fmt"
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
	_, _ = req.ToBytes()
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
	req.Client = sdk.HttpClient
	req.SetHeader("User-Agent", sdk.UserAgent)
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
	req.Client = sdk.HttpClient
	req.SetHeader("User-Agent", sdk.UserAgent)
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
	req.Client = sdk.HttpClient
	req.SetHeader("User-Agent", sdk.UserAgent)
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
	req.Client = sdk.HttpClient
	req.SetHeader("User-Agent", sdk.UserAgent)

	resp, err := req.Do()
	defer resp.Body.Close()
	if err == nil && resp.StatusCode == 200 {
		return true
	}
	return false
}

/* *********** * ************* */

// 获取服务器时间
func GetServerTime() int64 {
	req, err := dhttp.Get("https://a.jd.com//ajax/queryServerData.html")
	if err != nil {
		log.Panicln("GetServerTime", err.Error())
		return 0
	}
	type Ret struct {
		ServerTime int64 //1609878734768
	}
	var r Ret
	if err = req.ToJSON(&r); err != nil {
		log.Panicln("获取京东服务器时间失败", err.Error())
		return 0
	}
	return r.ServerTime
}

// 获取商品信息
func GetSeckillInitInfo(skuId, skuNum string) (*InitData, error) {
	log.Println("获取秒杀商品初始化信息...")
	req, err := dhttp.NewRequest("https://marathon.jd.com/seckillnew/orderService/pc/init.action", "POST")
	if err != nil {
		log.Panicln("GetSeckillInitInfo", err.Error())
		return nil, err
	}
	req.Client = sdk.HttpClient
	req.SetHeader("User-Agent", sdk.UserAgent)
	req.SetHeader("Host", "marathon.jd.com")

	req.WriteParam(url.Values{"sku": {skuId}, "num": {skuNum}, "isModifyAddress": {"false"}})

	var initData InitData
	if err = req.ToJSON(&initData); err != nil {
		log.Println("初始化秒杀信息失败", err.Error())
		return nil, err
	} else if len(initData.AddressList) == 0 {
		log.Println("初始化秒杀信息失败, AddressList为空")
		return nil, errors.New("初始化秒杀信息失败, AddressList为空")
	}
	return &initData, nil
}

// 获取秒杀链接
func GetKillUrl(skuId string) string {
	log.Println("获取秒杀商品链接...")
	req, err := dhttp.Get(dhttp.BuildURLParams("https://itemko.jd.com/itemShowBtn", url.Values{
		"skuId": {skuId}, "callback": {genCallback()}, "from": {"pc"}, "_": {genTime()},
	}))
	if err != nil {
		log.Println("获取秒杀商品链接失败", err.Error())
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
	if err = req.ToJSON(&r); err != nil {
		log.Println("获取秒杀商品链接失败, url为空")
		return ""
	}

	//https://divide.jd.com/user_routing?skuId=8654289&sn=c3f4ececd8461f0e4d7267e96a91e0e0&from=pc
	url := strings.ReplaceAll(r.Url, "divide", "marathon")
	//https://marathon.jd.com/captcha.html?skuId=8654289&sn=c3f4ececd8461f0e4d7267e96a91e0e0&from=pc
	url = strings.ReplaceAll(url, "user_routing", "captcha.html")
	log.Println("获取秒杀商品链接成功", url)
	return url
}

// 请求秒杀链接
func RequestKillUrl(skuId, killUrl string) {
	log.Println("请求秒杀商品链接...")
	req, err := dhttp.Get(killUrl)
	if err != nil {
		log.Println("请求秒杀商品链接失败", err.Error())
		return
	}
	req.Client = sdk.HttpClient
	req.SetHeader("User-Agent", sdk.UserAgent)
	req.SetHeader("Host", "marathon.jd.com")
	req.SetHeader("Referer", fmt.Sprintf("https://item.jd.com/%s.html", skuId))

	resp, err := req.Do()
	defer resp.Body.Close()
	if err == nil && resp.StatusCode == 200 {
		log.Println("请求秒杀商品链接成功")
	}
}

// 访问抢购订单结算页面
func SeckillPage(skuId, skuNum string) {
	log.Println("访问抢购订单结算页面...")
	req, err := dhttp.Get(dhttp.BuildURLParams("https://marathon.jd.com/seckill/seckill.action", url.Values{
		"sku": {skuId}, "num": {skuNum}, "rid": {genTime()},
	}))
	if err != nil {
		log.Println("访问抢购订单结算页面失败", err.Error())
		return
	}
	req.Client = sdk.HttpClient
	req.SetHeader("User-Agent", sdk.UserAgent)
	req.SetHeader("Host", "marathon.jd.com")
	req.SetHeader("Referer", fmt.Sprintf("https://item.jd.com/%s.html", skuId))

	resp, err := req.Do()
	defer resp.Body.Close()
	if err == nil && resp.StatusCode == 200 {
		log.Println("访问抢购订单结算页面成功")
	}

}

// 提交订单
func SubmitSeckillOrder(eid, fp, skuId, skuNum, pwd string, initData *InitData) bool {
	log.Println("提交商品订单...")
	req, err := dhttp.NewRequest(dhttp.BuildURLParams("https://marathon.jd.com/seckill/seckill.action", url.Values{"skuId": {skuId}}), "POST")
	if err != nil {
		log.Println("提交商品订单失败", err.Error())
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
		Success      bool //todo 类型等待验证
		ErrorMessage string
		OrderId      string
		ResultCode   string
		TotalMoney   string
		PcUrl        string
	}

	var r Ret
	if err = req.ToJSON(&r); err != nil {
		log.Println("提交商品订单失败2", err.Error())
		return false
	}
	if r.Success {
		log.Println(fmt.Sprintf("抢购成功，订单号:%s, 总价:%s, 电脑端付款链接:%s", r.OrderId, r.TotalMoney, r.PcUrl))
	}
	return r.Success
}
