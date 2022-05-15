package request

type RequestSendMiniProgramHB struct {

	ActName     string   `xml:"act_name"`
	MchBillNO   string   `xml:"mch_billno"`
	MchID       string   `xml:"mch_id"`
	NonceStr    string   `xml:"nonce_str"`
	NotifyWay   string   `xml:"notify_way"`
	ReOpenID    string   `xml:"re_openid"`
	Remark      string   `xml:"remark"`
	SendName    string   `xml:"send_name"`
	TotalAmount int   `xml:"total_amount"`
	TotalNum    int   `xml:"total_num"`
	Wishing     string   `xml:"wishing"`
	WXAppID     string   `xml:"wxappid"`
	Sign        string   `xml:"sign"`
}

