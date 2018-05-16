package wx

type JsCode2SessionResponse struct {
	Openid     string `json:"openid"`
	SessionKey string `json:"session_key"`
	Unionid    string `json:"unionid"`
	ErrorCode  int32  `json:"errorcode"`
	ErrMsg     string `json:"errmsg"`
}

type UnifiedOrderResponse struct {
	ReturnCode string `json:"return_code"`
	ReturnMsg  string `json:"return_msg"`
	DeviceInfo string `json:"device_info"`
	Appid      string `json:"appid"`
	MchID      string `json:"mch_id"`
	NonceStr   string `json:"nonce_str"`
	Sign       string `json:"sign"`
	ResultCode string `json:"result_code"`
	ErrCode    string `json:"err+_code"`
	ErrCodeDes string `json:"err_code_des"`
	TradeType  string `json:"trade_type"`
	PrePayID   string `json:"prepay_id"`
	CodeUrl    string `json:"code_url"`
}
