package pay

import (
	"encoding/json"
	"fmt"
	"time"
	"unicode/utf8"

	"github.com/yaotian/gowechat/mch/base"
	"github.com/yaotian/gowechat/util"
)

//OrderInput 下单
//官网文档 https://pay.weixin.qq.com/wiki/doc/api/jsapi.php?chapter=9_1
type OrderInput struct {
	OpenID      string //trade_type=JSAPI时（即公众号支付），此参数必传，此参数为微信用户在商户对应appid下的唯一标识
	Body        string //String(128)
	OutTradeNum string //String(32) 20150806125346 商户系统内部订单号，要求32个字符内，只能是数字、大小写字母_-|*@ ，且在同一个商户号下唯一。
	TotalFee    int    //分为单位
	IP          string
	NotifyURL   string //异步接收微信支付结果通知的回调地址，通知url必须为外网可访问的url，不能携带参数
	TradeType   string //JSAPI，NATIVE，APP
	ProductID   string //trade_type=NATIVE时（即扫码支付），此参数必传
}

//Config jssdk支付时需要的配置
type WxPayConfig struct {
	AppID     string `json:"appId"`
	TimeStamp string `json:"timeStamp"`
	NonceStr  string `json:"nonceStr"`
	Package   string `json:"package"`
	SignType  string `json:"signType"`
	PaySign   string `json:"paySign"`
	resultMap map[string]string
}

//ToString wx.chooseWXPay content
func (c *WxPayConfig) ToString() (str string) {
	return fmt.Sprintf(`
		timestamp: %s,
    nonceStr: '%s',
    package: '%s',
    signType: '%s',
    paySign: '%s', 
		`, c.TimeStamp, c.NonceStr, c.Package, c.SignType, c.PaySign)
}

//ToJSON WeixinJSBridge  json content
func (c *WxPayConfig) ToJSON() (str string) {
	js, err := json.Marshal(c)
	if err == nil {
		return string(js)
	}
	return
}

//ToMap result map[string]string
func (c *WxPayConfig) ToMap() (m map[string]string) {
	return c.resultMap
}

/*GetJsAPIConfig 前端JsAPI支付时,需要提交的信息
 */
func (c *Pay) GetJsAPIConfig(order OrderInput) (config *WxPayConfig, err error) {
	err = c.checkOrder(order)
	if err != nil {
		return
	}
	var prepayID string
	prepayID, err = c.getPrepayID(order)
	if err != nil {
		return
	}

	nocestr := util.RandomStr(8)
	timestamp := fmt.Sprint(time.Now().Unix())

	result := make(map[string]string)
	result["appId"] = c.AppID
	result["timeStamp"] = timestamp
	result["nonceStr"] = nocestr
	result["package"] = "prepay_id=" + prepayID
	result["signType"] = "MD5"

	sign := base.Sign(result, c.MchAPIKey, nil)
	result["paySign"] = sign

	config = new(WxPayConfig)
	config.NonceStr = util.RandomStr(8)
	config.TimeStamp = fmt.Sprint(time.Now().Unix())
	config.AppID = c.AppID
	config.Package = "prepay_id=" + prepayID
	config.SignType = "MD5"
	config.PaySign = sign

	config.resultMap = result

	return
}

// 调用 UnifiedOrder 获得 prepayID
func (c *Pay) getPrepayID(order OrderInput) (prepayID string, err error) {
	input := c.createUnifiedOrderMap(order)
	var result map[string]string
	if result, err = c.UnifiedOrder(input); err == nil { //有prepay_id
		prepayID := result["prepay_id"]
		if prepayID != "" {
			return prepayID, nil
		}
		err = fmt.Errorf("prepayID is empty")
	}
	return
}

func (c *Pay) createUnifiedOrderMap(order OrderInput) (input map[string]string) {
	input = make(map[string]string)
	input["appid"] = c.AppID               //设置微信分配的公众账号ID
	input["mch_id"] = c.MchID              //设置微信支付分配的商户号
	input["nonce_str"] = util.RandomStr(5) //设置随机字符串，不长于32位。推荐随机数生成算法
	input["body"] = order.Body             //获取商品或支付单简要描述的值

	input["out_trade_no"] = order.OutTradeNum       //设置商户系统内部的订单号,32个字符内、可包含字母, 其他说明见商户订单号
	input["total_fee"] = util.ToStr(order.TotalFee) //设置订单总金额，只能为整数，详见支付金额
	input["spbill_create_ip"] = order.IP            //设置APP和网页支付提交用户端ip，Native支付填调用微信支付API的机器IP。
	input["notify_url"] = order.NotifyURL           //设置接收微信支付异步通知回调地址

	input["trade_type"] = order.TradeType
	//设置取值如下：JSAPI，NATIVE，APP，详细说明见参数规定

	if order.ProductID != "" {
		input["product_id"] = order.ProductID //这个
	}

	input["openid"] = order.OpenID //设置trade_type=JSAPI，此参数必传，用户在商户appid下的唯一标识。下单前需要调用【网页授权获取用户信息】接口获取到用户的Openid

	//sign
	sign := base.Sign(input, c.MchAPIKey, nil)
	input["sign"] = sign
	return
}

func (c *Pay) checkOrder(order OrderInput) (err error) {
	tradeType := order.TradeType
	if tradeType != "JSAPI" && tradeType != "APP" && tradeType != "NATIVE" {
		return fmt.Errorf("tradeType is invalid")
	}
	if tradeType == "NATIVE" {
		if order.ProductID == "" {
			err = fmt.Errorf("Native TradeType need ProductID")
			return
		}
	}
	if tradeType == "JSAPI" {
		if order.OpenID == "" {
			err = fmt.Errorf("OpenID can not be empty when pay mode is JSAPI")
			return
		}
	}
	if utf8.RuneCountInString(order.Body) > 128 || order.Body == "" {
		err = fmt.Errorf("Body is invalid. Size can not exceed 128.")
		return
	}
	if utf8.RuneCountInString(order.OutTradeNum) > 32 || order.OutTradeNum == "" {
		err = fmt.Errorf("OutTradeNum is invalid. Size can not exceed 128.")
		return
	}
	if order.TotalFee <= 0 {
		err = fmt.Errorf("Order TotalFee is invalid.")
		return
	}
	if order.IP == "" {
		err = fmt.Errorf("Order IP is invalid.")
		return
	}
	if order.NotifyURL == "" {
		err = fmt.Errorf("Notify URL is invalid.")
		return
	}

	return

}
