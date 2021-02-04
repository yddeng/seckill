package jd

type InitData struct {
	AddressList  []Address
	InvoiceInfo  *InvoiceInfo
	SeckillSkuVO *SeckillSkuVO
	Token        string
}

type Address struct {
	Id            int
	Name          string
	ProvinceId    int
	CityId        int
	CountyId      int
	TownId        int
	AddressDetail string
	Email         string
	Mobile        string
	MobileKey     string
}

type InvoiceInfo struct {
	InvoiceTitle       int
	InvoiceContentType int
	InvoicePhone       string
	InvoicePhoneKey    string
}

type SeckillSkuVO struct {
	ExtMap ExtMap
}

type ExtMap struct {
	YuShou string
}
