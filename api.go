package wx

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"
)

const (
	LoginURL = "https://api.weixin.qq.com/sns/jscode2session"
	OrderURL = "https://api.mch.weixin.qq.com/pay/unifiedorder"
)

type APIProxy struct {
	config *Config
}

func NewAPIProxy(config *Config) *APIProxy {
	return &APIProxy{
		config: config,
	}
}

func (api *APIProxy) Login(jsCode string) (*JsCode2SessionResponse, error) {
	params := make(map[string]string)
	params["appid"] = api.config.Appid
	params["secret"] = api.config.Secret
	params["js_code"] = jsCode
	params["grant_type"] = "authorization_code"

	var response JsCode2SessionResponse
	err := api.request(LoginURL, params, response)
	return &response, err
}

func (api *APIProxy) UnifiedOrder(openid, tradeNo, body, totalFee, ipaddr string) (*PaymentRequest, error) {
	params := make(map[string]string)
	params["appid"] = api.config.Appid
	params["mch_id"] = api.config.MchID
	params["nonce_str"] = RandomString(32)
	params["body"] = body
	params["out_trade_no"] = tradeNo
	params["total_fee"] = totalFee
	params["spbill_create_ip"] = ipaddr
	params["notify_url"] = api.config.Notify
	params["trade_type"] = api.config.TradeType
	params["openid"] = openid
	params["sign_type"] = "MD5"
	params["sign"] = GenerateSign(api.config.Secret, params)

	var response UnifiedOrderResponse
	var request PaymentRequest

	err := api.request(OrderURL, params, response)
	if response.ReturnCode == "SUCCESS" &&
		response.ResultCode == "SUCCESS" &&
		len(response.PrePayID) > 0 {
		requsetParams := make(map[string]string)
		requsetParams["appId"] = api.config.Appid
		requsetParams["timeStamp"] = strconv.Itoa(int(time.Now().Unix()))
		requsetParams["nonceStr"] = response.NonceStr
		requsetParams["signType"] = "MD5"
		requsetParams["package"] = fmt.Sprintf("prepay_id=%s", response.PrePayID)
		paySign := GenerateSign(api.config.Secret, requsetParams)

		request = PaymentRequest{
			TimeStamp: requsetParams["timeStamp"],
			NonceStr:  requsetParams["nonceStr"],
			Package:   requsetParams["package"],
			SignType:  requsetParams["signType"],
			PaySign:   paySign,
		}
	}
	return &request, err
}

func (api *APIProxy) request(requestURL string, params map[string]string, result interface{}) error {
	req, err := http.NewRequest("GET", requestURL, nil)
	if err != nil {
		return err
	}

	q := req.URL.Query()
	for k, v := range params {
		q.Add(k, v)
	}

	req.URL.RawQuery = q.Encode()

	log.Println(req.URL.String())
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	json.Unmarshal([]byte(body), result)

	return nil
}
