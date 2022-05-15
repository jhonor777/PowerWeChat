package transfer

import (
	"crypto"
	"github.com/ArtisanCloud/PowerLibs/object"
	"github.com/ArtisanCloud/PowerLibs/security/sign"
	payment "github.com/ArtisanCloud/PowerWeChat/src/payment/kernel"
	"github.com/ArtisanCloud/PowerWeChat/src/payment/transfer/request"
	"github.com/ArtisanCloud/PowerWeChat/src/payment/transfer/response"
)

type Client struct {
	*payment.BaseClient
}

func NewClient(app *payment.ApplicationPaymentInterface) *Client {
	return &Client{
		payment.NewBaseClient(app),
	}
}

// Query MerchantPay to balance.
// https://pay.weixin.qq.com/wiki/doc/api/tools/mch_pay.php?chapter=14_3
func (comp *Client) QueryBalanceOrder(partnerTradeNo string) (*response.ResponseGetTransferInfo, error) {
	config := (*comp.App).GetConfig()

	result := &response.ResponseGetTransferInfo{}

	params := &object.StringMap{
		"appid":            config.GetString("app_id", ""),
		"mch_id":           config.GetString("mch_id", ""),
		"partner_trade_no": partnerTradeNo,
	}

	endpoint := comp.Wrap("/mmpaymkttransfers/gettransferinfo")
	_, err := comp.SafeRequest(endpoint, params, "POST", &object.StringMap{}, false, result)

	return result, err
}

// Send MerchantPay to balance.
// https://pay.weixin.qq.com/wiki/doc/api/tools/mch_pay.php?chapter=14_2
func (comp *Client) ToBalance(data *request.RequestTransferToBalance) (*response.ResponseTransferToBalance, error) {

	result := &response.ResponseTransferToBalance{}

	config := (*comp.App).GetConfig()

	//externalRequest := (*comp.App).GetExternalRequest()
	//if params.SpBillCreateIP == "" {
	//	params.SpBillCreateIP = externalRequest.Host
	//}

	params, err := object.StructToStringMap(data, "xml")
	if err != nil {
		return nil, err
	}
	base := &object.StringMap{
		"mchid":     config.GetString("mch_id", ""),
		"mch_appid": config.GetString("app_id", ""),
	}

	params = object.MergeStringMap(params, base)
	endpoint := comp.Wrap("mmpaymkttransfers/promotion/transfers")
	_, err = comp.SafeRequest(endpoint, params, "POST", &object.StringMap{}, nil, result)

	return result, err
}

// Query MerchantPay order to BankCard.
// https://pay.weixin.qq.com/wiki/doc/api/tools/mch_pay.php?chapter=14_3
func (comp *Client) QueryBankCardOrder(partnerTradeNo string) (*response.ResponseGetTransferInfo, error) {
	config := (*comp.App).GetConfig()

	result := &response.ResponseGetTransferInfo{}

	params := &object.StringMap{
		"mch_id":           config.GetString("mch_id", ""),
		"partner_trade_no": partnerTradeNo,
	}

	endpoint := comp.Wrap("/mmpaymkttransfers/query_bank")
	_, err := comp.SafeRequest(endpoint, params, "POST", &object.StringMap{}, false, result)

	return result, err
}

// Send MerchantPay to BankCard.
// https://pay.weixin.qq.com/wiki/doc/api/tools/mch_pay.php?chapter=24_2
func (comp *Client) ToBankCard(data *request.RequestToBankCard) (*response.ResponseTransfer, error) {

	result := &response.ResponseTransfer{}

	config := (*comp.App).GetConfig()

	rsaSigner, err := sign.NewRSASigner(crypto.SHA256)
	if err != nil {
		return nil, err
	}
	rsaSigner.RSAEncryptor.PublicKeyPath = config.GetString("rsa_public_key_path", "")
	_, err = rsaSigner.RSAEncryptor.LoadPublicKeyByPath()
	if err != nil {
		return nil, err
	}

	buffer, err := rsaSigner.RSAEncryptor.Encrypt([]byte(data.EncBankNO))
	if err != nil {
		return nil, err
	}
	data.EncBankNO = string(buffer)
	if err != nil {
		return nil, err
	}

	buffer, err = rsaSigner.RSAEncryptor.Encrypt([]byte(data.EncTrueName))
	if err != nil {
		return nil, err
	}
	data.EncTrueName = string(buffer)
	if err != nil {
		return nil, err
	}

	params, err := object.StructToStringMap(data, "xml")
	if err != nil {
		return nil, err
	}

	endpoint := comp.Wrap("mmpaysptrans/pay_bank")
	_, err = comp.SafeRequest(endpoint, params, "POST", &object.StringMap{}, nil, result)

	return result, err
}
