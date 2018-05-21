package wx

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"
)

const (
	LoginURL         = "https://api.weixin.qq.com/sns/jscode2session"
	OrderURL         = "https://api.mch.weixin.qq.com/pay/unifiedorder"
	GetTokenURL      = "https://api.weixin.qq.com/cgi-bin/token"
	SetUserStorgeURL = "https://api.weixin.qq.com/wxa/set_user_storage"
)

var (
	NoSupportMethod   = errors.New("nonsupport method")
	UnifiedOrderError = errors.New("unified order error")
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
	err := api.request(http.MethodGet, "", LoginURL, params, response)
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

	err := api.request(http.MethodGet, "", OrderURL, params, response)
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
	} else {
		err = UnifiedOrderError
	}
	return &request, err
}

func (api APIProxy) GetToken() (*GetTokenResponse, error) {
	params := make(map[string]string)
	params["appid"] = api.config.Appid
	params["secret"] = api.config.Secret
	params["grant_type"] = "client_credential"

	var response GetTokenResponse
	err := api.request(http.MethodGet, "", GetTokenURL, params, response)
	return &response, err
}

func (api APIProxy) SetUserStorge(openid, accessToken, sessionKey string, kvList string) (*SetUserStorgeResponse, error) {
	params := make(map[string]string)
	params["appid"] = api.config.Appid
	params["access_token"] = accessToken
	params["signature"] = GenerateLoginStatusSign(kvList, sessionKey)
	params["sig_method"] = "hmac_sha256"

	var response SetUserStorgeResponse

	err := api.request(http.MethodPost, kvList, SetUserStorgeURL, params, response)
	return &response, err
}

func (api *APIProxy) request(method, requestBody, requestURL string, params map[string]string, result interface{}) error {
	var req *http.Request
	var err error
	if method == http.MethodGet {
		req, err = http.NewRequest(http.MethodGet, requestURL, nil)
	} else if method == http.MethodPost {
		requsetBody := bytes.NewReader([]byte(requestBody))
		req, err = http.NewRequest(http.MethodGet, requestURL, requsetBody)
	} else {
		return NoSupportMethod
	}

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
